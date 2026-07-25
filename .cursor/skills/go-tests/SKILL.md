---
name: go-tests
description: >-
  Test wat CLI, SDK, codec, adapter, merge, and filesystem behavior using the
  repository's current architecture. Use when adding tests, reviewing test
  quality, or implementing hook/CLI behavior.
---

# Go tests

Follow [CONTRIBUTING.md](../../../CONTRIBUTING.md) and the extension checklists
in [docs/architecture.md](../../../docs/architecture.md).

## Two layers

| Layer | Location | Rule |
|---|---|---|
| Unit | Package `*_test.go` | Inject I/O; `t.TempDir()`; inline JSON; synthetic manifests. No init + replace + real hook build. |
| E2E | [`e2e/`](../../../e2e/) | Build real `cmd/wat`; scaffold `.wat/`; assert CLI streams, exit codes, and fixture hook runs. |

## General rules

- Assert observable behavior; no placeholder or tautological tests.
- Prefer table-driven cases for protocol matrices and flag validation.
- Use descriptive `t.Run` names.
- Use `t.Helper()` only inside helpers.
- Use `t.TempDir()` and injected dependencies for filesystem/process behavior.
- Store reusable native JSON for e2e under `e2e/testdata/`; keep small focused
  payloads inline. Project `wat test` cases live under `.wat/testdata/`.

## SDK and protocol changes

For a result-producing native event, cover as applicable:

- event-name peeking and typed decode;
- missing/invalid/unknown event behavior;
- builder output and fluent `With*` fields;
- encoded JSON and process exit code;
- output merge and terminal `Stop`;
- nil handler registration.

For portable behavior, test all three adapters and assert the expanded native
events/handler counts from `run.Inspect`.

When a combine function accepts structs containing maps or slices, snapshot the
caller's input and assert it remains unchanged as well as checking the result.

## CLI changes

- Keep command help on stdout with exit 0.
- Keep invalid usage on stderr with exit 1.
- Test internal behavior with injected I/O, environment, filesystem, and
  process functions in the owning package.
- Cover the public command and hook contract in `e2e/` (build and run `cmd/wat`).
- Assert distinctive output, the correct stream, and the exact exit code.

## Verification

```bash
go test ./path/to/touched/package
go test ./e2e
go build ./cmd/wat
go vet ./...
go test ./...
golangci-lint run ./...
```
