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

## Caller immutability

Functions that take a struct **by value** and return a combined result must not mutate reference-type fields owned by the caller (maps, slices, pointers).

Applies to merge/combine helpers such as `Merge(a, b Result) Result`: `a` is copied, but `a.Env` still aliases the caller's map unless the implementation clones first.

When adding or reviewing these APIs, assert **both**:

1. **Output** — merged fields have the expected values.
2. **Input unchanged** — snapshot caller-owned maps/slices before the call and compare after.

```go
t.Run("env merge", func(t *testing.T) {
    a := Result{Env: map[string]string{"A": "1", "B": "2"}}
    b := Result{Env: map[string]string{"B": "override", "C": "3"}}
    origA := maps.Clone(a.Env)

    got := Merge(a, b)

    if got.Env["B"] != "override" {
        t.Fatalf("merged env = %v", got.Env)
    }
    if !maps.Equal(a.Env, origA) {
        t.Fatalf("caller env mutated: got %v, want %v", a.Env, origA)
    }
})
```

Use `maps.Clone` / `maps.Equal` for maps and `slices.Clone` / `slices.Equal` for slices. For pointers, compare dereferenced values or identity as appropriate to the contract.

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
