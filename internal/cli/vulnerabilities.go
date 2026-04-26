package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/urugus/baseline-cli/internal/baseline"
)

func newVulnerabilitiesCommand(global *globalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vulnerabilities",
		Aliases: []string{"vulns", "vuln"},
		Short:   "Read vulnerabilities",
	}

	cmd.AddCommand(newVulnerabilitiesListCommand(global))
	cmd.AddCommand(newVulnerabilitiesGetCommand(global))
	return cmd
}

func newVulnerabilitiesListCommand(global *globalOptions) *cobra.Command {
	var page int
	var perPage int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List vulnerabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := baseline.NewClient(baseline.ClientOptions{BaseURL: global.baseURL})
			if err != nil {
				return err
			}

			response, raw, err := client.ListVulnerabilities(context.Background(), baseline.ListVulnerabilitiesOptions{
				Page:    page,
				PerPage: perPage,
			})
			if err != nil {
				return err
			}
			if global.json {
				return printJSON(raw)
			}
			return printVulnerabilityList(response)
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Items per page")
	return cmd
}

func newVulnerabilitiesGetCommand(global *globalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a vulnerability",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := baseline.NewClient(baseline.ClientOptions{BaseURL: global.baseURL})
			if err != nil {
				return err
			}

			response, raw, err := client.GetVulnerability(context.Background(), args[0])
			if err != nil {
				return err
			}
			if global.json {
				return printJSON(raw)
			}
			return printVulnerabilityDetail(response.Data)
		},
	}
	return cmd
}

func printJSON(raw []byte) error {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(encoded))
	return nil
}

func printVulnerabilityList(response baseline.PageResponse[baseline.Vulnerability]) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NB\tSEVERITY\tSTATUS\tPROJECT\tASSET\tTITLE\tID")
	for _, vuln := range response.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			vuln.NB,
			vuln.Severity,
			vuln.Status,
			vuln.Project.Name,
			vuln.Asset.DisplayName,
			vuln.Title,
			vuln.ID,
		)
	}
	fmt.Fprintf(w, "\npage=%d perPage=%d total=%d\n", response.Page, response.PerPage, response.Total)
	return w.Flush()
}

func printVulnerabilityDetail(vuln baseline.Vulnerability) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", vuln.ID)
	fmt.Fprintf(w, "NB:\t%s\n", vuln.NB)
	fmt.Fprintf(w, "Severity:\t%s\n", vuln.Severity)
	fmt.Fprintf(w, "Status:\t%s\n", vuln.Status)
	fmt.Fprintf(w, "Project:\t%s (%s)\n", vuln.Project.Name, vuln.Project.ID)
	fmt.Fprintf(w, "Asset:\t%s (%s)\n", vuln.Asset.DisplayName, vuln.Asset.ID)
	fmt.Fprintf(w, "Title:\t%s\n", vuln.Title)
	if vuln.TitleJP != "" {
		fmt.Fprintf(w, "Title JP:\t%s\n", vuln.TitleJP)
	}
	fmt.Fprintf(w, "Custom ID:\t%s\n", vuln.CustomID)
	fmt.Fprintf(w, "Created At:\t%s\n", vuln.CreatedAt)
	fmt.Fprintf(w, "Updated At:\t%s\n", vuln.UpdatedAt)
	if len(vuln.CVEs) > 0 {
		fmt.Fprintf(w, "CVEs:\t%v\n", vuln.CVEs)
	}
	if len(vuln.Refs) > 0 {
		fmt.Fprintf(w, "Refs:\t%v\n", vuln.Refs)
	}
	return w.Flush()
}
