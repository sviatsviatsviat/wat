# wat

**wat** is a Go CLI and library stack for authoring hook scripts that run under Claude Code, GitHub Copilot, and Cursor. The redesign provides a unified `agenthooks` library, per-agent SDKs, and a `wat` command that compiles and installs user hook scripts.

## Build

```bash
go build ./cmd/wat
```

On Windows the toolchain writes `wat.exe`; on Unix, `wat`. Avoid `go build -o wat` on Windows — use `./cmd/wat` or `-o wat.exe`.

## Test and lint

```bash
go vet ./...
go test ./...
golangci-lint run ./...
```

Install golangci-lint once:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
```

Browse package docs locally:

```bash
go doc github.com/sviatsviatsviat/wat/agenthooks
```

## Design docs

Architecture and implementation tasks live in a local-only `plan/` directory (gitignored). They are not part of a standard clone.

## License

Apache License 2.0 — see [LICENSE](LICENSE).
