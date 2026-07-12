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
- Helpers call `t.Helper()`.
- Fixtures live in `testdata/` (JSON payloads, golden files).

## Hook and CLI tests

- Prefer building and exercising `cmd/wat` over calling unexported helpers when validating end-to-end behavior.
- Inject I/O (`io.Reader` / `io.Writer`, fake getenv) instead of globals.
- No machine-specific paths unless hermetic with `t.TempDir()`.
- For stdout/stderr protocols (hook JSON, help text), assert against test-local expected strings or golden files — not duplicated production constants unless intentional.

## Running tests

```bash
go test ./...
go test -race ./...   # optional locally
```

## Exit criteria

Tests should fail when behavior regresses and pass when behavior is correct. Avoid tests that only assert mocks call mocks.
