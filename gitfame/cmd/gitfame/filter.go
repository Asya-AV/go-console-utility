//go:build !solution

package main

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func matchAnyExtension(file string) bool {
	fileExtension := strings.ToLower(filepath.Ext(file))
	for _, ext := range FlagExtensions {
		if strings.ToLower(ext) == fileExtension {
			return true
		}
	}
	return false
}

func matchAnyLanguage(file string, langExtensions map[string][]string) bool {
	fileExtension := strings.ToLower(filepath.Ext(file))
	for _, lang := range FlagLanguages {
		if ext, ok := langExtensions[strings.ToLower(lang)]; ok {
			for _, e := range ext {
				if e == fileExtension {
					return true
				}
			}
			continue
		}
	}
	return false
}

func matchAnyPattern(file string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, file); matched {
			return true
		}
	}
	return false
}

func FilterFiles(files []string) ([]string, error) {
	langExtensions := viper.GetStringMapStringSlice("languages")
	var filteredFiles []string

	for _, file := range files {
		if len(FlagExtensions) > 0 && !matchAnyExtension(file) {
			continue
		}

		if len(FlagLanguages) > 0 && !matchAnyLanguage(file, langExtensions) {
			continue
		}

		if len(FlagExclude) > 0 && matchAnyPattern(file, FlagExclude) {
			continue
		}

		if len(FlagRestrictTo) > 0 && !matchAnyPattern(file, FlagRestrictTo) {
			continue
		}

		filteredFiles = append(filteredFiles, file)
	}

	return filteredFiles, nil
}
