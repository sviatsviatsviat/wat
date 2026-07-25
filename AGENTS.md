# Repository instructions

`wat` is a Go 1.26 module (`github.com/sviatsviatsviat/wat`) for authoring and
installing hook scripts across Claude Code, GitHub Copilot, and Cursor.

Read these committed references before changing behavior:

- [Architecture](docs/architecture.md) for package ownership, dependencies, and
  extension patterns.
- [SDK API](docs/sdk.md) for the public author-facing contract.
- [Using wat](docs/usage.md) for CLI behavior.
- [Agent protocols](docs/agent-formats.md) for codec and normalization facts.
- [Contributing](CONTRIBUTING.md) for tests, docs, changelog, commits, and CI.

## Package map

- `cmd/wat`: thin CLI entrypoints; behavior lives in focused
  `cmd/wat/internal/*` packages.
- `sdk/run`: process dispatch and registration manifests.
- `sdk/claude`, `sdk/copilot`, `sdk/cursor`: native protocol owners and public
  facades.
- `sdk/agnostic`: portable model and fan-out adapters; depends on all native
  SDKs.
- `sdk/agnostic/tools`: canonical portable tool names and typed inputs.
- `internal/hookkit`: module-private codec, handler, merge, and normalization
  machinery.
- `e2e`: public CLI and hook end-to-end tests (real `wat` binary + `.wat/`
  project). Package-local tests stay injected and non-social.
- `testdata/fixtures`: native hook payloads consumed by package tests and `e2e`.
- `.wat/testdata`: project-local fixtures and optional `*.expect.json` sidecars
  for `wat test` (excluded from the hook binary cache key).

## Non-negotiable boundaries

- Native SDKs must not import `sdk/agnostic`.
- Native wire decoding, output encoding, merge, and exit behavior belong to the
  owning native SDK.
- Portable APIs expose only behavior all three agents can represent.
- CLI install/test/doctor derive registered events from `run.Inspect`; do not
  maintain a second event table.
- Keep registration and package initialization free of external side effects.
- CLI command files parse flags and map exit codes; inject I/O, filesystem,
  environment, and process dependencies in internal implementations.
- Every exported identifier has godoc. Package overviews go in `doc.go`.
- Result/merge operations must not mutate caller-owned maps or slices.

## Required verification

```bash
go vet ./...
go test ./...
golangci-lint run ./...
go build ./cmd/wat
```

Install the pinned linter outside the Dev Container:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
```

Use focused package tests while iterating, then run the complete suite
(including `./e2e`). New public behavior requires real behavior tests and
documentation in the same change. The Linux Dev Container is the authoritative
full-suite environment; some CLI and e2e tests rely on Unix utilities and are
not Windows-hermetic.

## Change hygiene

- User-visible CLI, SDK, hook protocol, or security behavior belongs in
  `CHANGELOG.md`. Internal refactors, CI, tests, docs-only edits, and agent
  tooling do not.
- Update README/guides/godoc together when their described behavior changes.
- Keep GitHub Actions and tool versions exactly pinned. A dependency bump
  updates every reference to that dependency in the same change.
- Use conventional commits with imperative subjects: `feat:`, `fix:`,
  `refactor:`, `test:`, `docs:`, or `chore:`.
- Pull request summaries cover the full branch and list verification.

Local design notes may use the gitignored `plan/` directory, but current
architecture and behavior must be documented in committed files.
