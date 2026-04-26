package baseline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadEnvValue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".zshenv.local")
	if err := os.WriteFile(path, []byte(`
# unrelated
export BASELINE_API_KEY="api_test"
`), 0o600); err != nil {
		t.Fatal(err)
	}

	value, err := readEnvValue(path, "BASELINE_API_KEY")
	if err != nil {
		t.Fatal(err)
	}
	if value != "api_test" {
		t.Fatalf("expected api_test, got %q", value)
	}
}

func TestReadEnvValueWithoutExport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".zshenv.local")
	if err := os.WriteFile(path, []byte(`BASELINE_API_KEY='api_test' # comment`), 0o600); err != nil {
		t.Fatal(err)
	}

	value, err := readEnvValue(path, "BASELINE_API_KEY")
	if err != nil {
		t.Fatal(err)
	}
	if value != "api_test" {
		t.Fatalf("expected api_test, got %q", value)
	}
}
