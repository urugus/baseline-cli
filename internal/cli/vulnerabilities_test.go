package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/urugus/baseline-cli/internal/baseline"
)

func TestOutputFormat(t *testing.T) {
	tests := []struct {
		name    string
		opts    globalOptions
		want    string
		wantErr bool
	}{
		{
			name: "defaults to table",
			want: "table",
		},
		{
			name: "accepts json",
			opts: globalOptions{format: "json"},
			want: "json",
		},
		{
			name: "normalizes ndjson",
			opts: globalOptions{format: " NDJSON "},
			want: "ndjson",
		},
		{
			name: "json flag wins",
			opts: globalOptions{json: true, format: "table"},
			want: "json",
		},
		{
			name:    "rejects unsupported format",
			opts:    globalOptions{format: "yaml"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := outputFormat(&tt.opts)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("format = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrintVulnerabilityListNDJSON(t *testing.T) {
	response := vulnerabilityListResponse{
		PageResponse: baseline.PageResponse[baseline.Vulnerability]{
			Data: []baseline.Vulnerability{
				{
					ID:       "vuln-1",
					NB:       "NB-1",
					Severity: "critical",
					Status:   "open",
					Project:  baseline.ProjectRef{Name: "project-a"},
					Asset:    baseline.AssetRef{DisplayName: "asset-a"},
					Title:    "Example vulnerability",
				},
				{
					ID:       "vuln-2",
					NB:       "NB-2",
					Severity: "high",
					Status:   "fixed",
					Project:  baseline.ProjectRef{Name: "project-b"},
					Asset:    baseline.AssetRef{DisplayName: "asset-b"},
					Title:    "Another vulnerability",
				},
			},
			Page:    1,
			PerPage: 2,
			Total:   10,
		},
		ServerTotal: 20,
	}

	var out bytes.Buffer
	if err := printVulnerabilityListNDJSON(&out, response); err != nil {
		t.Fatalf("printVulnerabilityListNDJSON returned error: %v", err)
	}

	got := out.String()
	wantLines := []string{
		`{"id":"vuln-1","nb":"NB-1","severity":"critical","status":"open","project":"project-a","asset":"asset-a","title":"Example vulnerability"}`,
		`{"id":"vuln-2","nb":"NB-2","severity":"high","status":"fixed","project":"project-b","asset":"asset-b","title":"Another vulnerability"}`,
	}
	if strings.TrimSpace(got) != strings.Join(wantLines, "\n") {
		t.Fatalf("output = %q, want %q", strings.TrimSpace(got), strings.Join(wantLines, "\n"))
	}
}
