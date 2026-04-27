---
name: prepare-release
description: Prepare a release PR for baseline-cli by updating CHANGELOG.md, running tests, and opening a PR against the base branch. Does NOT create tags or publish releases.
disable-model-invocation: true
allowed-tools: Bash, Read, Write, Edit, Glob, Grep
arguments: [version, previous_tag, base_ref]
argument-hint: <version> "<previous_tag>" <base_ref>
---

You are preparing a release PR for the baseline-cli repository.

## Hard boundaries
- Do not create or push git tags.
- Do not publish a GitHub Release.
- Do not update the Homebrew tap.
- Do not change runtime behavior unless required to update release metadata.

## Inputs
- `$version` — Proposed version (e.g. `v0.3.0`)
- `$previous_tag` — Previous tag (may be empty; if empty, use the latest local tag matching `v*`)
- `$base_ref` — Base ref (branch to target the PR against, e.g. `main`)

## Steps

### 1. Determine the comparison range
- If `$previous_tag` is non-empty, use it as `range_start`.
- Otherwise run `git tag --list 'v*' --sort=-v:refname | head -n 1` and use that.
- Abort with a clear error if no previous tag can be determined.

### 2. Inspect commits in `${range_start}..HEAD`
Use `git log --format='%H%x09%s%x09%b'` (or similar) to get subjects and bodies. Read enough context to write good notes — do not just dump subject lines.

### 3. Create or update `CHANGELOG.md`
Insert a new section for `$version` at the top (immediately after the `# Changelog` heading if present; otherwise create the file with that heading). Preserve all prior sections.

**Format of the new section:**

```
## $version — YYYY-MM-DD

### Added
- ...

### Changed
- ...

### Fixed
- ...

### Removed
- (only if applicable)

Operational follow-up: run the Release workflow manually with $version after this PR is merged.
```

Use today's UTC date.

### 4. Release-notes quality bar
Write notes a user reading the GitHub release page will understand without context. Specifically:

- **Group by intent**, not by commit prefix mechanics. Use `Added` / `Changed` / `Fixed` / `Removed` sections (Keep a Changelog convention). Map `feat:` → Added, `fix:` → Fixed, refactors/perf/internal config → Changed, deletions → Removed.
- **Skip noise**: omit `chore:`, `docs:`, `test:`, `ci:`, `style:`, dependency bumps with no user-visible effect, and merge commits — unless they affect users (e.g. a doc change that announces a deprecation).
- **User-facing language**: rewrite from the user's perspective. Bad: `Refactor config loader`. Good: `Improve config loading speed on large repos`. Bad: `Add filter flag`. Good: `Filter vulnerabilities by severity with --severity`.
- **Highlight breaking changes**: if any commit body contains `BREAKING CHANGE:` or removes/renames a flag/command/config key, add a `### Breaking changes` section at the very top of the version's section, before `Added`. State what broke and how to migrate.
- **One bullet per change**, not per commit. Squash duplicate or follow-up commits ("fix typo in previous commit") into a single bullet.
- **No commit hashes or PR numbers in bullets** — the CHANGELOG is for users, not git archaeology. (The PR body can list raw commits if needed.)

### 5. Run tests
Run `go test ./...`. Abort if anything fails.

### 6. Commit, push, and open the PR

Compute the branch name by stripping a leading `v` from `$version` (e.g. `v0.3.0` → `release/prepare-0.3.0`):

```bash
version="$version"
branch="release/prepare-${version#v}"
git checkout -B "$branch"
git config user.name  "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git add CHANGELOG.md
git commit -m "Prepare release $version"
git push --force-with-lease --set-upstream origin "$branch"
```

Then open the PR with `gh pr create`:

- `--base $base_ref --head "$branch"`
- `--title "Prepare release $version"`
- `--body` containing:
  - **Summary** — 1–3 bullets describing the user-visible contents of the release (mirror the CHANGELOG highlights, not raw commits).
  - **Validation** — `go test ./...`
  - A final paragraph: `Tag and GitHub Release publication remain a human-triggered step through the Release workflow.`

If a PR with the title `Prepare release $version` already exists open against `$base_ref`, do not create a duplicate — update the existing branch instead and report the existing PR number.
