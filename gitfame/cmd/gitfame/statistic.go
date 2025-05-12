//go:build !solution

package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Statistic struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func isHex(s string) bool {
	for _, c := range s {
		if !(('0' <= c && c <= '9') || ('a' <= c && c <= 'f')) {
			return false
		}
	}
	return true
}

func runGitBlame(file string) ([]byte, error) {
	cmd := exec.Command("git", "-C", FlagRepository, "blame", "--line-porcelain", FlagRevision, "--", file)
	return cmd.Output()
}

func runGitLog(file string) ([]byte, error) {
	cmd := exec.Command("git", "-C", FlagRepository, "log", "--pretty=fuller", FlagRevision, "--", file)
	return cmd.Output()
}

func processBlameOutput(lines []string, stats map[string]map[string]int) bool {
	var currentCommit, currentAuthor string
	var commitLines int
	isEmpty := true

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) != "" {
			isEmpty = false
		}

		if len(line) >= 40 && isHex(line[:40]) {
			currentCommit = line[:40]
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				commitLines, _ = strconv.Atoi(parts[3])
			} else {
				commitLines = 0
			}
			continue
		}

		if currentCommit == "" {
			continue
		}

		if !FlagUseCommitter && strings.HasPrefix(line, "author ") {
			currentAuthor = extractName(line, "author ")
			updateStats(stats, currentAuthor, currentCommit, commitLines)
			currentCommit = ""
			continue
		}
		if FlagUseCommitter && strings.HasPrefix(line, "committer ") {
			currentAuthor = extractName(line, "committer ")
			updateStats(stats, currentAuthor, currentCommit, commitLines)
			currentCommit = ""
		}
	}

	return isEmpty
}

func extractName(line, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}

func updateStats(stats map[string]map[string]int, author, commit string, lines int) {
	if author == "" {
		return
	}
	if _, ok := stats[author]; !ok {
		stats[author] = make(map[string]int)
	}
	stats[author][commit] += lines
}

func processEmptyFileLog(lines []string, stats map[string]map[string]int) {
	var author, commit string
	for _, line := range lines {
		if strings.HasPrefix(line, "commit ") {
			commit = strings.TrimSpace(strings.TrimPrefix(line, "commit "))
		}
		if !FlagUseCommitter && strings.HasPrefix(line, "Author: ") {
			author = cleanAuthorName(strings.TrimPrefix(line, "Author: "))
		}
		if FlagUseCommitter && strings.HasPrefix(line, "Commit: ") {
			author = cleanAuthorName(strings.TrimPrefix(line, "Commit: "))
		}

		if author != "" && commit != "" {
			if _, ok := stats[author]; !ok {
				stats[author] = make(map[string]int)
			}
			stats[author][commit] = 0
			break
		}
	}
}

func cleanAuthorName(s string) string {
	return regexp.MustCompile(`\s*<[^>]+>\s*`).ReplaceAllString(strings.TrimSpace(s), "")
}

func getFileStats(file string) (map[string]map[string]int, error) {
	stats := make(map[string]map[string]int)

	output, err := runGitBlame(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	isEmpty := processBlameOutput(lines, stats)

	if !isEmpty {
		return stats, nil
	}

	output, err = runGitLog(file)
	if err != nil {
		return nil, err
	}

	lines = strings.Split(string(output), "\n")
	processEmptyFileLog(lines, stats)

	return stats, nil
}

func DoStatistic(files []string) []Statistic {
	stats := make(map[string]*Statistic)
	commitTracker := make(map[string]map[string]struct{})
	fileTracker := make(map[string]map[string]struct{})

	for _, file := range files {
		fileStats, err := getFileStats(file)
		if err != nil {
			return nil
		}

		for author, commitMap := range fileStats {
			if _, ok := stats[author]; !ok {
				stats[author] = &Statistic{Name: author}
				commitTracker[author] = make(map[string]struct{})
				fileTracker[author] = make(map[string]struct{})
			}

			for commitHash, count := range commitMap {
				stats[author].Lines += count

				if _, exists := commitTracker[author][commitHash]; !exists {
					commitTracker[author][commitHash] = struct{}{}
					stats[author].Commits++
				}
			}

			if _, exists := fileTracker[author][file]; !exists {
				fileTracker[author][file] = struct{}{}
				stats[author].Files++
			}
		}
	}

	var result []Statistic
	for _, stat := range stats {
		result = append(result, *stat)
	}
	return result
}
