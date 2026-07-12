---
name: go-tests
description: >-
  Write Go tests for wat following repo conventions. Use when adding tests,
  reviewing test quality, or implementing hook/CLI behavior.
---

# Go tests

## Principles

- Assert **real behavior**. No empty, placeholder, or tautological tests.
- New exported API or behavior change ships with tests in the same change.

## Structure

- Table-driven tests for decode/encode matrices and flag parsing:

```go
tests := []struct {
    name string
    // fields...
}{
    {name: "deny destructive shell"},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // ...
    })
}
```

- Test names: `Test<Function>_<scenario>` with subtests via `t.Run`.
- **`t.Helper()` only in helper functions** (`buildWat`, `moduleRoot`, fixture loaders) — not in top-level `TestXxx` or `t.Run` callbacks.
- Fixtures live in `testdata/` (JSON payloads, golden files).

## Hook and CLI tests

- Prefer building and exercising `cmd/wat` over calling unexported helpers when validating end-to-end behavior.
- Inject I/O (`io.Reader` / `io.Writer`, fake getenv) instead of globals.
- No machine-specific paths unless hermetic with `t.TempDir()`.
- For stdout/stderr protocols (hook JSON, help text), assert against test-local expected strings or golden files — not duplicated production constants unless intentional.

### CLI table template

Use per-case `wantOutput` with distinctive substrings — not generic checks like `strings.Contains(out, "wat")`:

```go
tests := []struct {
    name       string
    args       []string
    wantErr    bool
    wantCode   int
    wantOutput string
}{
    {
        name:       "help",
        args:       []string{"help"},
        wantOutput: "Run with -h or help for this message.",
    },
    {
        name:       "unknown_command",
        args:       []string{"nosuchcommand"},
        wantErr:    true,
        wantCode:   1,
        wantOutput: "Commands will be added in later releases",
    },
}
```

### CLI stream conventions

Production `cmd/wat` behavior:

- **Help** (`-h`, `--help`, `help`) → write usage to **stdout**, exit 0
- **Invalid usage** (no args, unknown command) → write usage to **stderr**, exit 1

`exec.Command(...).CombinedOutput()` merges streams for tests; assert on distinctive usage lines via `wantOutput`.

## Running tests

```bash
go test ./...
go test -race ./...   # optional locally
go vet ./... && go test ./...
```

## Exit criteria

Tests should fail when behavior regresses and pass when behavior is correct. Avoid tests that only assert mocks call mocks.
