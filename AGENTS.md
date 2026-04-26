# Repository Guidelines

## Project Structure & Module Organization

This repository contains a small Go CLI for read-only IssueHunt Baseline access. The executable entry point is `cmd/baseline/main.go`. CLI commands live in `internal/cli`, API access and response types live in `internal/baseline`, and local configuration helpers live in `internal/config`. Tests are colocated with the package they cover, using Go’s `*_test.go` convention. Build outputs are written to `baseline` or `dist/`; do not commit generated binaries.

## Build, Test, and Development Commands

- `make build`: builds a local `./baseline` binary with version metadata.
- `./baseline version`: quick smoke test for the built binary.
- `make test`: runs `go test ./...` across all packages.
- `make fmt`: runs `gofmt -w .` on the repository.
- `make install`: installs the CLI from `./cmd/baseline`.
- `make release VERSION=v0.1.0`: builds cross-platform release artifacts and checksums under `dist/`.

For local API use, set `BASELINE_API_KEY` or configure it with `baseline config set api-key`.

## Coding Style & Naming Conventions

Use standard Go formatting with tabs as produced by `gofmt`; run `make fmt` before committing. Keep package names short and lowercase. Export only APIs needed across packages. Cobra command constructors follow the existing `newXCommand` pattern, and tests should use descriptive names such as `TestAPIKeyPrefersEnvironment`.

## Testing Guidelines

Use the standard Go `testing` package. Prefer focused package-level tests next to the code under test, for example `internal/config/config_test.go`. Use `t.TempDir()` and `t.Setenv()` for config and environment isolation. Run `make test` before opening a pull request. Add tests for new flags, config behavior, API response parsing, and error handling paths.

## Commit & Pull Request Guidelines

Recent commits use short, imperative subject lines, often followed by a PR number after merge, for example `Add vulnerability list filters (#2)` or `Fix prepare release Claude workflow (#6)`. Keep commit subjects specific and under roughly 72 characters.

Pull requests should include a concise summary, testing performed, and any release or configuration impact. Link related issues when applicable. For CLI behavior changes, include example commands or output. Screenshots are usually unnecessary unless GitHub workflow or documentation rendering is affected.

## Security & Configuration Tips

This CLI is intentionally read-only: preserve the GET-only safety model and avoid adding arbitrary URL access. Never print API keys. Keep credentials in `BASELINE_API_KEY` or the local config file, and ensure config-file changes maintain restrictive permissions.
