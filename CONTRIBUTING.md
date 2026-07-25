# Contributing to wat

Changes are expected to keep code, tests, godoc, user documentation, and the
changelog aligned. Start with the package boundary in
[Architecture](docs/architecture.md), then use this guide for the development
workflow.

## Development environment

The module requires Go 1.26.

Use the pinned Linux [Dev Container](.devcontainer/README.md) for the
authoritative full verification suite. The Dev Container stores the workspace
in the `wat-workspace` Docker volume at `/workspaces/wat`; it is not a Windows
bind mount.

A local Go environment is useful for editing and focused package tests, but
some CLI tests execute Unix utilities such as `echo` and `sleep` and therefore
are not currently Windows-hermetic.

Install the pinned linter when working outside the container:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
```

Before submitting a change:

```bash
go vet ./...
go test ./...
golangci-lint run ./...
go build ./cmd/wat
```

CI also runs `govulncheck`.

## Find the right package

| Change | Primary location |
|---|---|
| CLI parsing/help | `cmd/wat` |
| CLI behavior | focused package under `cmd/wat/internal` |
| Hook process dispatch/manifest | `sdk/run` |
| Claude/Copilot/Cursor protocol | owning `sdk/<agent>` |
| Cross-agent behavior | `sdk/agnostic` and all three adapters |
| Shared private protocol machinery | `internal/hookkit` |

Do not put native wire behavior in `sdk/agnostic`, and do not make a native SDK
depend on `sdk/agnostic`. Installation must use `run.Inspect` as the source of
registered native events.

## Implementing changes

### Public SDK changes

- Follow the vertical-slice patterns in [Architecture](docs/architecture.md).
- Add godoc for every exported identifier.
- Use hook-scoped result builders; do not expose freely constructible wire
  output structs.
- Preserve native input while normalizing portable fields.
- Add behavior tests for decode, builder, merge, encode, exit code, and adapter
  mapping as applicable.
- Update [SDK API](docs/sdk.md), [Agent protocols](docs/agent-formats.md), and
  the changelog in the same change.

### CLI changes

- Keep `cmd/wat/<command>.go` limited to flags, help, and exit-code mapping.
- Inject I/O, environment, filesystem, and process functions into internal
  implementations.
- Help goes to stdout with exit 0. Invalid usage goes to stderr with exit 1.
- Use `t.TempDir()` for filesystem tests and build the real CLI for public
  end-to-end behavior.
- Update command help, [Using wat](docs/usage.md), and the changelog.

### Internal changes

Internal refactors still need tests when behavior or invariants can regress.
They do not need a changelog entry unless users or hook authors observe a
change.

## Tests

Tests must assert behavior rather than implementation bookkeeping.

- Prefer table-driven tests for protocol matrices and flag validation.
- Use subtests with descriptive scenario names.
- Put reusable native JSON under `testdata/`; keep narrowly focused payloads
  inline when easier to understand.
- Use `t.Helper()` only in helper functions.
- Assert both the combined output and caller input immutability for merge
  helpers that receive maps or slices by value.
- Include the expected stdout/stderr stream and process exit code in CLI tests.
- For portable changes, test each agent adapter and the native registration
  expansion returned by `run.Inspect`.

Run a focused package test while iterating, then the complete verification
suite before handoff.

## Documentation and godoc

User-facing behavior has three documentation layers:

1. command/package godoc for exact API behavior;
2. guides under `docs/` for workflows and concepts;
3. `CHANGELOG.md` for release-level outcomes.

Update all affected layers in the same pull request. Avoid documenting planned
types, events, packages, or commands as if they exist. Examples should compile
or be close enough to copy into `.wat/hooks.go` without hidden setup.

Godoc comments are complete sentences that begin with the exported identifier.
Package overviews belong in `doc.go`. CI enforces exported comments with
`revive`.

## Changelog policy

The changelog is for behavior available to end users and SDK consumers:

- CLI capabilities, flags, discovery, installation, cache, and exit behavior;
- public SDK events, builders, and compatibility;
- native protocol support and behavior-affecting fixes;
- security fixes.

Do not add entries for repository scaffolding, CI, lint configuration, tests,
agent instructions, documentation-only edits, or internal refactors with no
visible effect. Entries should describe what a user can do and the behavior
they can rely on, not the internal packages used to implement it.

Use Keep a Changelog sections. Until there is a published baseline, describe
the accumulated initial release under `Added`; use `Changed` or `Removed` only
relative to behavior that users could obtain from a published release.

## Version pinning

Tool and GitHub Action versions are exact. Never introduce `@latest`, branch
tags, or floating major tags.

When bumping a dependency, update every reference to that same dependency in
the same pull request. For example, a golangci-lint bump must update the
workflow, Dev Container, this guide, and agent instructions if they cite it.

Audit workflow action tags with:

```bash
rg 'uses:\s*[\w./-]+@v\d+$' .github/workflows
```

The command must produce no matches.

## Commits and pull requests

Use an imperative conventional commit subject:

```text
feat(agnostic): add portable lifecycle hook
fix(wat): preserve unrelated installed handlers
docs: rewrite SDK architecture guide
```

Supported prefixes are `feat`, `fix`, `refactor`, `test`, `docs`, and `chore`.
Use the body to explain why the change is needed when the subject is not enough.

A pull request should:

- summarize the full branch, not only the last commit;
- call out user-visible behavior and compatibility;
- list verification performed;
- include documentation and changelog updates when required;
- keep unrelated cleanup out of the change.
