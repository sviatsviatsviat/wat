---
name: task-planning
description: >-
  Plan wat changes against the current architecture with explicit scope,
  compatibility, tests, documentation, and completion checks. Use when creating
  or updating plan/tasks/*.md files or drafting implementation plans.
---

# Task planning

Plans may live under the gitignored `plan/` directory when a durable local note
is useful. Do not assume a particular design document or numbered task series
exists.

## Before planning

Read the relevant current references:

- [Architecture](../../../docs/architecture.md) for ownership and invariants;
- [SDK API](../../../docs/sdk.md) for public contract;
- [Using wat](../../../docs/usage.md) for CLI behavior;
- [Agent protocols](../../../docs/agent-formats.md) for native mappings;
- [Contributing](../../../CONTRIBUTING.md) for completion requirements.

Inspect the implementation and tests. Documentation describes the intended
contract, while code and tests establish what currently ships; call out any
disagreement explicitly.

## Plan structure

Include:

1. goal and user-visible outcome;
2. in-scope and explicitly deferred work;
3. affected packages and dependency direction;
4. ordered implementation steps;
5. compatibility or native-protocol differences;
6. tests by layer;
7. documentation/changelog impact;
8. definition of done.

## Definition of done

- [ ] Package boundaries in `docs/architecture.md` remain valid
- [ ] Registerable hook packages keep `RegisterHandler` in `bind.go` (never in
      `event.go`); observe-only still has `bind.go`
- [ ] Hook event decode avoids custom `UnmarshalJSON` for duration/presence
- [ ] Public API has godoc and behavior tests
- [ ] Portable behavior is tested for all three adapters and manifest expansion
- [ ] CLI help, streams, and exit codes are tested when affected
- [ ] Caller-owned maps/slices remain unchanged by merge/combine APIs
- [ ] README/guides/protocol docs are updated where behavior changed
- [ ] Changelog contains only user-visible outcomes
- [ ] Exact dependency pins remain synchronized
- [ ] `go vet ./...`, `go test ./...`, `golangci-lint run ./...`, and
      `go build ./cmd/wat` pass

Tailor the list to the task; do not add speculative packages or future events
to committed documentation.
