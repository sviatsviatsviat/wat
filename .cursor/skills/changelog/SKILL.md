---
name: changelog
description: >-
  Update CHANGELOG.md using Keep a Changelog rules. Use when adding release
  notes, documenting shipped features, or editing the [Unreleased] section.
---

# Changelog

Follow [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) in [CHANGELOG.md](../../CHANGELOG.md).

## Sections

Use **Added**, **Changed**, **Deprecated**, **Removed**, **Fixed**, **Security** under `[Unreleased]` or a version heading.

## Rules

- **Added** — new work in the current release cycle.
- **Changed** / **Removed** — only when comparing to an **already published** release.
- While `[Unreleased]` has no prior public release, put corrections in **Added**, not **Changed**.
- Write affirmative, reader-focused entries. Describe what exists and what users can do.
- Do not use negative framing or list features that were never shipped.

## Example

```markdown
### Added

- `agenthooks` unified Event model with dialect detection.
```

Not: "Removed experimental X" when X never shipped publicly.
