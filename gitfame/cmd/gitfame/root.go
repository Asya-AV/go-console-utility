//go:build !solution

package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	FlagRepository, FlagRevision, FlagOrderBy, FlagFormat      string
	FlagUseCommitter                                           bool
	FlagExtensions, FlagLanguages, FlagExclude, FlagRestrictTo []string
)

var rootCmd = &cobra.Command{
	Use:   "gitfame",
	Short: "Summarizing and printing collaborators, based on the number of contributions in git repository",
	Long:  "Summarizing and printing collaborators in git repository with number of their lines, commits and number of their files",
	RunE:  runGitfame,
}

func checkValid() {
	validOrderBy := map[string]bool{
		"lines":   true,
		"commits": true,
		"files":   true,
	}

	if !validOrderBy[FlagOrderBy] {
		log.Printf("invalid value for --order-by: %v. Please, use valid values: lines, commits, files", FlagOrderBy)
		os.Exit(1)
	}
}

func runGitfame(cmd *cobra.Command, args []string) error {
	checkValid()
	files, errGet := GetFiles()
	if errGet != nil {
		return errGet
	}

	filteredFiles, errFilter := FilterFiles(files)
	if errFilter != nil {
		return errFilter
	}

	stats := DoStatistic(filteredFiles)

	SortStats(stats)

	return PrintStats(stats, os.Stdout)
}

func initConfig() {
	file, err := os.Open("../../configs/language_extensions.json")
	if err != nil {
		log.Printf("Warning: config file not found: %v. Using empty configuration.", err)
		return
	}
	defer file.Close()

	var languages []struct {
		Name       string   `json:"name"`
		Type       string   `json:"type"`
		Extensions []string `json:"extensions"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&languages); err != nil {
		log.Printf("Warning: invalid config format: %v. Using empty configuration.", err)
		return
	}

	extensionsMap := make(map[string][]string)
	for _, lang := range languages {
		extensionsMap[strings.ToLower(lang.Name)] = lang.Extensions
	}

	viper.Set("languages", extensionsMap)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVar(&FlagRepository, "repository", ".", "path to Git repository")
	rootCmd.Flags().StringVar(&FlagRevision, "revision", "HEAD", "pointer to Git commit")
	rootCmd.Flags().StringVar(&FlagOrderBy, "order-by", "lines", "sort key (lines, commits, files)")
	rootCmd.Flags().BoolVar(&FlagUseCommitter, "use-committer", false, "use committer instead of author")
	rootCmd.Flags().StringVar(&FlagFormat, "format", "tabular", "output format (tabular, csv, json, json-lines)")
	rootCmd.Flags().StringSliceVar(&FlagExtensions, "extensions", nil, "comma-separated list of file extensions")
	rootCmd.Flags().StringSliceVar(&FlagLanguages, "languages", nil, "comma-separated list of languages")
	rootCmd.Flags().StringSliceVar(&FlagExclude, "exclude", nil, "comma-separated list of exclude Glob patterns")
	rootCmd.Flags().StringSliceVar(&FlagRestrictTo, "restrict-to", nil, "comma-separated list of include Glob patterns")

	cobra.OnInitialize(initConfig)
}
