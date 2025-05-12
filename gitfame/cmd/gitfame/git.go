//go:build !solution

package main

import (
	"os/exec"
	"strings"
)

func GetFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", FlagRevision)
	cmd.Dir = FlagRepository
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}
