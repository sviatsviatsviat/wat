---
name: conventional-commits
description: >-
  Draft intentional conventional commits for wat changes. Use when committing,
  summarizing changes for commit, or the user asks for commit message help.
---

# Conventional commits

Use:

```text
<type>(<optional-scope>): <imperative subject>

<optional body explaining why>
```

Supported types:

| Type | Use |
|---|---|
| `feat` | New user-visible CLI or SDK behavior |
| `fix` | Correction to user-visible or internal behavior |
| `refactor` | Code change without behavior change |
| `test` | Tests only |
| `docs` | Documentation only |
| `chore` | Build, CI, tooling, or dependency maintenance |

Useful scopes include `wat`, `run`, `agnostic`, `claude`, `copilot`, and
`cursor`.

Use imperative mood, omit the trailing period, and keep the subject concise.
Use the body for motivation, compatibility, or a non-obvious tradeoff.

Examples:

```text
feat(agnostic): add portable pre-compact hook
fix(wat): preserve unrelated installed handlers
docs: explain native SDK registration patterns
```
