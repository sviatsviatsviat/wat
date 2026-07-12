## Context

**wat** is a Go monorepo (`github.com/sviatsviatsviat/wat`, Go 1.26) for hook scripts across Claude Code, GitHub Copilot, and Cursor. This is a greenfield redesign: the old `internal/*` Cursor-only layout is not preserved.

**Packages:**

- **`cmd/wat`** — CLI entrypoint (`init`, `install`, `run`, `port`, `test`, `doctor` in later tasks).
- **`agenthooks`** — Aggregated library: unified `Event`/`Result`, dialect codecs, runner, and `portconfig`. Stdlib only.
- **`claudehook`**, **`copilothook`**, **`cursorhook`** — Per-agent SDKs; independent packages usable without `agenthooks`.
- **`internal/cmdast`** — Shell AST helpers (task 22; not present yet).

**Module graph:** `wat` depends on `agenthooks`. Per-agent SDKs do not depend on `agenthooks`.

**Design docs:** Local-only `plan/` directory (gitignored). Prototype code in `plan/agenthooks/` is reference material for tasks 02–10.

**Build and CI:**

```bash
go build ./cmd/wat
go vet ./...
go test ./...
golangci-lint run ./...
```

Install golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

Browse API docs: `go doc github.com/sviatsviatsviat/wat/agenthooks`

## Rules

- **Layering:** Keep `agenthooks` stdlib-only. Per-agent SDKs stay independent. Put CLI wiring in `cmd/wat`; do not reintroduce the old `internal/core` + host packages layout unless a later task specifies it.
- **Godoc:** Every exported identifier needs a godoc comment. Package overviews belong in `doc.go`. See `.cursor/rules/godoc.mdc` and the godoc skill.
- **Tests:** Assert real behavior; no placeholder tests. See `.cursor/skills/go-tests/SKILL.md` and `.cursor/rules/go-tests.mdc`.
- **Changelog:** Keep a Changelog format. **Added** for new work; **Changed** / **Removed** only vs a published release. Affirmative, reader-focused entries. See the changelog skill.
- **Commits:** Conventional prefixes (`feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`). Imperative subject; body explains *why*. See the conventional-commits skill.
- **Pull requests:** Summarize the full branch, not only the latest commit.
- **Naming:** Clear identifiers (`args`, `programArgs`) over jargon (`argv`, `d`, `rest`) unless Unix argv is explicitly the topic.
