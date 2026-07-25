---
name: changelog
description: >-
  Update CHANGELOG.md using Keep a Changelog rules. Use when adding release
  notes, documenting shipped user-facing features, or editing the [Unreleased]
  section.
---

# Changelog

Follow [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) in [CHANGELOG.md](../../CHANGELOG.md).

## What belongs

Document **user-facing functionality** — behavior users or hook authors can rely on in released artifacts:

- CLI commands, flags, exit codes
- Public library API (`sdk/agnostic`, per-agent SDKs under `sdk/`)
- Hook protocol support, breaking API changes
- Security fixes affecting shipped behavior

Write affirmative, reader-focused entries. Describe what exists and what users can do.

## What to omit

Do **not** add changelog entries for internal or non-artifact work, including:

- Repo layout, module scaffolding, package stubs
- CI, lint, test harness, or agent tooling (`.cursor/`, skills, rules)
- Refactors, chores, or docs that do not change shipped behavior
- Negative framing or features that were never released

If `[Unreleased]` has no user-facing changes yet, leave the section empty (no placeholder bullets).

## Sections

Use **Added**, **Changed**, **Deprecated**, **Removed**, **Fixed**, **Security** under `[Unreleased]` or a version heading.

- **Added** — new user-facing capability in the current release cycle.
- **Changed** / **Removed** — only when comparing to an **already published** release.
- While `[Unreleased]` has no prior public release, put corrections in **Added**, not **Changed**.

## Example

```markdown
### Added

- `wat run` executes `.wat/hooks.go` against stdin hook JSON from any supported agent.
```

Not:

```markdown
### Added

- Monorepo scaffold with CI and golangci-lint.
```
