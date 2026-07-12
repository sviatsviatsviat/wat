---
name: godoc
description: >-
  Write godoc comments for exported Go API. Use when adding or changing public
  funcs, types, constants, vars, methods, or package doc.go files.
---

# Godoc

Every exported identifier in this repo must have a godoc comment. CI enforces this via golangci-lint revive `exported`.

## Package docs

Put a package overview in `doc.go`:

```go
// Package agenthooks provides a unified hook event model for coding agents.
package agenthooks
```

Use `doc.go` when the overview helps; small internal packages may omit it if nothing is exported yet.

## Symbol comments

- Start with the name of the thing being described.
- Complete sentence, proper punctuation.
- Avoid stutter: `// KindPreTool is ...` not `// KindPreTool KindPreTool is ...`

```go
// KindPreTool is the normalized category for pre-tool hook events.
KindPreTool Kind = "PreTool"
```

## Browse locally

```bash
go doc github.com/sviatsviatsviat/wat/agenthooks
go doc github.com/sviatsviatsviat/wat/agenthooks Kind
```

No doc server is required; use standard `go doc`.

## When adding exports

Ship godoc in the **same change** as the new exported symbol. Do not leave exports undocumented and rely on a follow-up.
