---
name: godoc
description: >-
  Document exported Go APIs and package contracts for wat's portable, native,
  runtime, and CLI-support packages. Use when adding or changing public funcs,
  types, constants, vars, methods, or package doc.go files.
---

# Godoc

Every exported function, type, constant, variable, and method needs a complete
godoc sentence whose **first word/phrase is the declaration name**.

```go
// Foo returns the canonical name.   // good (funcs/methods/vars/consts)
// This function returns the name.   // bad — lead with Foo
// The Foo type holds config.        // OK for types (article + name)
// A Foo is a config handle.         // OK for types (optional article + name)
// The FooMethod does work.          // bad on methods — no leading article
```

Do not start with filler such as `This`, `For`, `Returns`, or `Provides`
without the identifier. Types may use an optional leading `A`/`An`/`The`
when the type name follows immediately; funcs, methods, vars, and consts
must begin with the declaration name (no article). Later sentences in the
same comment block may start freely; only the opening sentence is
constrained.

CI enforces this with golangci-lint:

- revive `exported` — missing exported comments and lead-name form
  (`comment on exported … should be of the form "Name …"`)
- staticcheck `ST1020` (funcs/methods), `ST1021` (types), `ST1022` (vars/consts)

Unexported helpers and ordinary test callbacks do not need docstrings for CI.
Run `golangci-lint run ./...` locally so lead-name mistakes fail before
CodeRabbit. Ignore CodeRabbit's private/unexported docstring coverage
percentage; that metric is not a merge gate.

Put meaningful package contracts in `doc.go`. For public SDK packages, explain:

- how handlers are registered;
- result-producing versus observe-only signatures;
- how outputs are constructed and what nil means;
- whether the package is portable, native, or a process boundary.

For native hook packages under `sdk/<agent>/internal/hooks/`, keep layout
comments consistent with [docs/architecture.md](../../../docs/architecture.md):
`RegisterHandler` lives in `bind.go` for both observe-only and result-emitting
events; `event.go` is types + codec register/decode only.

Document native limitations and side effects where callers choose behavior,
such as output exit codes or environment-file writes. Do not expose internal
layout details in symbol comments unless callers need them.

Use examples for copyable author workflows and keep them consistent with
[docs/sdk.md](../../../docs/sdk.md).

Browse the rendered result:

```bash
go doc github.com/sviatsviatsviat/wat/sdk/agnostic
go doc github.com/sviatsviatsviat/wat/sdk/agnostic PreToolEvent
go doc github.com/sviatsviatsviat/wat/sdk/run
```

Ship godoc in the same change as the exported API.
