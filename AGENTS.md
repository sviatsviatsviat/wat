## Context

**wat** is a Go monorepo (`github.com/sviatsviatsviat/wat`, Go 1.26) for hook scripts across Claude Code, GitHub Copilot, and Cursor. This is a greenfield redesign: the old `internal/*` Cursor-only layout is not preserved.

**Packages:**

- **`cmd/wat`** — CLI entrypoint (`init`, `install`, `run`, `port`, `test`, `doctor` in later tasks).
- **`sdk/agnostic`** — Aggregated library: unified `Event`/`Result`, unexported inbound native→unified mapping beside each `On*` hook, and `On*` registration that fans out onto per-agent SDK Chains. Depends on `sdk/claude`, `sdk/copilot`, and `sdk/cursor`. Kind/event registries live under `sdk/agnostic/claude`, `copilot`, and `cursor`.
- **`cmd/wat/internal/portconfig`** — CLI-only hook config porting for `wat port` (not a public library API).
- **`sdk/claude`**, **`sdk/copilot`**, **`sdk/cursor`** — Per-agent SDKs; independent packages usable without `sdk/agnostic`.
- **`internal/cmdast`** — Shell AST helpers (task 22; not present yet).

**Module graph:** `wat` depends on `sdk/agnostic`. Per-agent SDKs do not depend on `sdk/agnostic`.

**Design docs:** Local-only `plan/` directory (gitignored) for task and design notes.

**Development environment:** Linux Dev Container — see [.devcontainer/README.md](.devcontainer/README.md). Workspace lives on named volume `wat-workspace` at `/workspaces/wat` (no Windows bind mount).

**Build and CI:**

```bash
go build ./cmd/wat
go vet ./...
go test ./...
golangci-lint run ./...
```

Verify locally: `go vet ./... && go test ./... && golangci-lint run ./...`

Install golangci-lint: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2`

Browse API docs: `go doc github.com/sviatsviatsviat/wat/sdk/agnostic`

**Documentation:** User-facing behavior belongs in [docs/](docs/) and [CHANGELOG.md](CHANGELOG.md). Update those when API or CLI behavior changes in the same PR as the code.

**Dependabot:** [`.github/dependabot.yml`](.github/dependabot.yml) opens weekly PRs for GitHub Actions, Go modules, Dev Container Features, and the devcontainer base image (`.devcontainer/Dockerfile`). Keep exact pinned versions. When bumping a dependency, update every reference to **that same dependency** in the same PR (e.g. all `golangci-lint` pins: CI `version:`, README, AGENTS; all `govulncheck@` pins in workflow files). Unrelated pins do not need to move with every Action bump unless compatibility requires it.

## Rules

- **Layering:** Per-agent SDKs stay independent of `sdk/agnostic` (stdlib only). `sdk/agnostic` may depend on the agent SDKs and registers portable handlers onto their Chains. Put CLI wiring in `cmd/wat`; do not reintroduce the old `internal/core` + host packages layout unless a later task specifies it.
- **Godoc:** Every exported identifier needs a godoc comment. Package overviews belong in `doc.go`. See `.cursor/rules/godoc.mdc` and the godoc skill.
- **Tests:** Assert real behavior; no placeholder tests. See `.cursor/skills/go-tests/SKILL.md` and `.cursor/rules/go-tests.mdc`.
- **Changelog:** User-facing functionality only (CLI, public API, hook behavior) — not scaffolding, CI, lint, or agent harness. Keep a Changelog format; **Added** for new capability; **Changed** / **Removed** only vs a published release. See the changelog skill.
- **Version pinning:** Pin exact tool and action versions in workflows and docs; no `@latest`. Bump all references together. See **Version pinning** in `.cursor/rules/wat-core.mdc` and `.cursor/skills/ci-pins/SKILL.md`.
- **Commits:** Conventional prefixes (`feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`). Imperative subject; body explains *why*. See the conventional-commits skill.
- **Pull requests:** Summarize the full branch, not only the latest commit.
- **Naming:** Clear identifiers (`args`, `programArgs`) over jargon (`argv`, `d`, `rest`) unless Unix argv is explicitly the topic.
