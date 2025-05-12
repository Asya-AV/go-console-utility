//go:build !solution

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

func printTabular(stats []Statistic, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	fmt.Fprintln(tw, "Name\tLines\tCommits\tFiles")
	for _, stat := range stats {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\n", stat.Name, stat.Lines, stat.Commits, stat.Files)
	}
	return tw.Flush()
}

func printCSV(stats []Statistic, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"Name", "Lines", "Commits", "Files"}); err != nil {
		return err
	}
	for _, stat := range stats {
		if err := cw.Write([]string{
			stat.Name,
			fmt.Sprint(stat.Lines),
			fmt.Sprint(stat.Commits),
			fmt.Sprint(stat.Files),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return nil
}

func printJSON(stats []Statistic, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(stats)
}

func printJSONLines(stats []Statistic, w io.Writer) error {
	enc := json.NewEncoder(w)
	for _, stat := range stats {
		if err := enc.Encode(stat); err != nil {
			return err
		}
	}
	return nil
}

func PrintStats(stats []Statistic, w io.Writer) error {
	switch FlagFormat {
	case "tabular":
		return printTabular(stats, os.Stdout)

	case "csv":
		return printCSV(stats, os.Stdout)

	case "json":
		return printJSON(stats, os.Stdout)

	case "json-lines":
		return printJSONLines(stats, os.Stdout)

	default:
		return fmt.Errorf("unknown format: %s", FlagFormat)
	}
}
