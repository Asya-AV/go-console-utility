//go:build !solution

package main

import "sort"

func SortStats(stats []Statistic) {
	switch FlagOrderBy {
	case "commits":
		sort.Slice(stats, func(i, j int) bool {
			if stats[i].Commits != stats[j].Commits {
				return stats[i].Commits > stats[j].Commits
			}
			if stats[i].Lines != stats[j].Lines {
				return stats[i].Lines > stats[j].Lines
			}
			if stats[i].Files != stats[j].Files {
				return stats[i].Files > stats[j].Files
			}
			return stats[i].Name < stats[j].Name
		})
	case "files":
		sort.Slice(stats, func(i, j int) bool {
			if stats[i].Files != stats[j].Files {
				return stats[i].Files > stats[j].Files
			}
			if stats[i].Lines != stats[j].Lines {
				return stats[i].Lines > stats[j].Lines
			}
			if stats[i].Commits != stats[j].Commits {
				return stats[i].Commits > stats[j].Commits
			}
			return stats[i].Name < stats[j].Name
		})
	default:
		sort.Slice(stats, func(i, j int) bool {
			if stats[i].Lines != stats[j].Lines {
				return stats[i].Lines > stats[j].Lines
			}
			if stats[i].Commits != stats[j].Commits {
				return stats[i].Commits > stats[j].Commits
			}
			if stats[i].Files != stats[j].Files {
				return stats[i].Files > stats[j].Files
			}
			return stats[i].Name < stats[j].Name
		})
	}
}
