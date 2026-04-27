package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type globalOptions struct {
	baseURL string
	json    bool
	format  string
}

func NewRootCommand() *cobra.Command {
	opts := &globalOptions{}

	cmd := &cobra.Command{
		Use:   "baseline",
		Short: "Read-only CLI for IssueHunt Baseline",
		Long:  "Read-only CLI for IssueHunt Baseline. This tool only implements GET-based commands.",
	}

	cmd.PersistentFlags().StringVar(&opts.baseURL, "base-url", "", "Baseline API base URL")
	cmd.PersistentFlags().BoolVar(&opts.json, "json", false, "Print raw JSON response")
	cmd.PersistentFlags().StringVar(&opts.format, "format", "table", "Output format: table, json, ndjson")

	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newVulnerabilitiesCommand(opts))
	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "baseline %s commit=%s date=%s\n", version, commit, date)
		},
	}
}
