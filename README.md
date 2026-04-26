# baseline-cli

Read-only CLI for IssueHunt Baseline.

This tool is intentionally scoped to reference operations. It only implements
GET-based API calls and does not expose arbitrary URL access.

## Install

With Homebrew:

```sh
brew tap urugus/tap
brew install baseline
```

From source:

```sh
go install github.com/urugus/baseline-cli/cmd/baseline@latest
```

For local development:

```sh
git clone https://github.com/urugus/baseline-cli.git
cd baseline-cli
make install
```

Or build a local binary:

```sh
make build
./baseline version
```

## Authentication

Set `BASELINE_API_KEY`.

```sh
export BASELINE_API_KEY='...'
```

If your shell stores it in `~/.zshenv.local`, source it before running this CLI
from automation:

```sh
source ~/.zshenv.local
```

## Usage

```sh
baseline vulnerabilities list
baseline vulnerabilities list --page 1 --per-page 20
baseline vulnerabilities list --asset TOKIUM/drwallet-worker --severity critical --all
baseline vulnerabilities list --asset-id 140e0253-9bbc-4f60-9ad2-e2742aa11b2a --severity critical --all
baseline vulnerabilities list --json
baseline vulnerabilities get <id>
baseline vulnerabilities get <id> --json
```

## Safety Policy

- Only GET requests are implemented.
- API paths are fixed by subcommand.
- Arbitrary URLs are not accepted.
- API keys are read from `BASELINE_API_KEY` and never printed.
- Create, update, and delete operations are intentionally out of scope.

## Release

Build release artifacts:

```sh
make release VERSION=v0.1.0
```

This creates platform binaries and `checksums.txt` under `dist/`.

To publish a GitHub Release, push a version tag:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds binaries for macOS, Linux, and Windows, then uploads
them with SHA-256 checksums.

CI and release jobs use GitHub Actions `ubuntu-slim` runners because the build is
lightweight and does not require Docker, service containers, or privileged
operations.

Maintainers can also run the release workflow manually with a version such as
`v0.3.0`. Manual releases create the tag, publish GitHub Release assets, and can
open a Homebrew tap update PR. The tap PR requires a repository secret named
`RELEASE_PAT` with push access to `urugus/homebrew-tap`.
