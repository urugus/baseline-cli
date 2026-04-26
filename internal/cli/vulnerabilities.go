package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	var all bool
	var severity string
	var status string
	var asset string
	var assetID string
	var project string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List vulnerabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			if all && cmd.Flags().Changed("page") {
				return fmt.Errorf("--page cannot be used with --all")
			}
			client, err := baseline.NewClient(baseline.ClientOptions{BaseURL: global.baseURL})
			if err != nil {
				return err
			}

			filter := vulnerabilityFilter{
				severity: severity,
				status:   status,
				asset:    asset,
				project:  project,
			}
			response, err := loadVulnerabilities(context.Background(), client, listOptions{
				Page:    page,
				PerPage: perPage,
				All:     all,
				AssetID: assetID,
				Filter:  filter,
			})
			if err != nil {
				return err
			}
			if global.json {
				return printVulnerabilityListJSON(response)
			}
			return printVulnerabilityList(response)
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Items per page")
	cmd.Flags().BoolVar(&all, "all", false, "Fetch all pages")
	cmd.Flags().StringVar(&severity, "severity", "", "Filter by severity")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&asset, "asset", "", "Filter by asset display name")
	cmd.Flags().StringVar(&assetID, "asset-id", "", "Filter by asset UUID on the server side")
	cmd.Flags().StringVar(&project, "project", "", "Filter by project name")
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

type listOptions struct {
	Page    int
	PerPage int
	All     bool
	AssetID string
	Filter  vulnerabilityFilter
}

type vulnerabilityFilter struct {
	severity string
	status   string
	asset    string
	project  string
}

func loadVulnerabilities(ctx context.Context, client *baseline.Client, opts listOptions) (baseline.PageResponse[baseline.Vulnerability], error) {
	if opts.PerPage <= 0 {
		opts.PerPage = 20
	}
	if !opts.All {
		response, _, err := client.ListVulnerabilities(ctx, baseline.ListVulnerabilitiesOptions{
			Page:    opts.Page,
			PerPage: opts.PerPage,
			AssetID: opts.AssetID,
		})
		if err != nil {
			return baseline.PageResponse[baseline.Vulnerability]{}, err
		}
		response.Data = filterVulnerabilities(response.Data, opts.Filter)
		return response, nil
	}

	var all []baseline.Vulnerability
	page := 1
	total := 0
	for {
		response, _, err := client.ListVulnerabilities(ctx, baseline.ListVulnerabilitiesOptions{
			Page:    page,
			PerPage: opts.PerPage,
			AssetID: opts.AssetID,
		})
		if err != nil {
			return baseline.PageResponse[baseline.Vulnerability]{}, err
		}
		total = response.Total
		all = append(all, filterVulnerabilities(response.Data, opts.Filter)...)
		if response.Page*response.PerPage >= response.Total || len(response.Data) == 0 {
			break
		}
		page++
	}

	return baseline.PageResponse[baseline.Vulnerability]{
		Data:    all,
		Page:    1,
		PerPage: len(all),
		Total:   total,
	}, nil
}

func filterVulnerabilities(vulnerabilities []baseline.Vulnerability, filter vulnerabilityFilter) []baseline.Vulnerability {
	if filter.isZero() {
		return vulnerabilities
	}
	filtered := make([]baseline.Vulnerability, 0, len(vulnerabilities))
	for _, vuln := range vulnerabilities {
		if filter.matches(vuln) {
			filtered = append(filtered, vuln)
		}
	}
	return filtered
}

func (f vulnerabilityFilter) isZero() bool {
	return f.severity == "" && f.status == "" && f.asset == "" && f.project == ""
}

func (f vulnerabilityFilter) matches(vuln baseline.Vulnerability) bool {
	if f.severity != "" && !equalFold(vuln.Severity, f.severity) {
		return false
	}
	if f.status != "" && !equalFold(vuln.Status, f.status) {
		return false
	}
	if f.asset != "" && !containsFold(vuln.Asset.DisplayName, f.asset) {
		return false
	}
	if f.project != "" && !containsFold(vuln.Project.Name, f.project) {
		return false
	}
	return true
}

func equalFold(value string, expected string) bool {
	return strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(expected))
}

func containsFold(value string, expected string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(strings.TrimSpace(expected)))
}

func printVulnerabilityListJSON(response baseline.PageResponse[baseline.Vulnerability]) error {
	encoded, err := json.MarshalIndent(response, "", "  ")
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
	fmt.Fprintf(w, "\nshown=%d total=%d\n", len(response.Data), response.Total)
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
