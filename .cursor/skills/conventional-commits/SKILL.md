---
name: conventional-commits
description: >-
  Draft git commit messages with conventional prefixes. Use when committing,
  summarizing changes for commit, or the user asks for commit message help.
---

# Conventional commits

## Format

```
<type>(<optional-scope>): <imperative subject>

<optional body explaining why>
```

## Types

| Prefix | Use for |
|--------|---------|
| `feat:` | New user-facing behavior or API |
| `fix:` | Bug fix |
| `refactor:` | Code change without behavior change |
| `test:` | Tests only |
| `docs:` | Documentation only |
| `chore:` | Build, CI, tooling, deps |

## Scope (optional)

Examples: `feat(wat):`, `feat(agnostic):`, `fix(cursor):`

## Subject

- Imperative mood: "add handler" not "added handler"
- No trailing period
- ~72 characters or less when practical

## Body

Explain **why** the change is needed. Reference issues or design sections when helpful.

## Examples

```
feat(agnostic): add unified Event Kind constants

Normalize event names across Claude, Copilot, and Cursor so one handler
registers once in Mux.On.
```

```
chore: wire golangci-lint and revive exported rule in CI
```
