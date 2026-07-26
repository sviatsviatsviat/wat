---
name: godoc
description: >-
  Document exported Go APIs and package contracts for wat's portable, native,
  runtime, and CLI-support packages. Use when adding or changing public funcs,
  types, constants, vars, methods, or package doc.go files.
---

# Godoc

Every exported function, type, constant, variable, and method needs a complete
godoc sentence beginning with its name. CI enforces this with `revive`
(`exported`) and with an 80% exported-godoc coverage gate:

```bash
go run ./tools/check-godoc-coverage -threshold=80 -list-missing
```

The coverage metric matches project rules: exported funcs, methods, types,
consts, and vars with AST doc comments. `_test.go`, `testdata/`, and generated
files are excluded; struct fields are not counted. Keep coverage ≥80% on every
PR—prefer adding missing exported godoc over lowering the threshold or gaming
excludes. CodeRabbit's PR "docstring coverage" warning is advisory (diff-scoped
and may include unexported symbols); this tool is authoritative in CI (workflow
step `Godoc coverage`).

Put meaningful package contracts in `doc.go`. For public SDK packages, explain:

- how handlers are registered;
- result-producing versus observe-only signatures;
- how outputs are constructed and what nil means;
- whether the package is portable, native, or a process boundary.

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
