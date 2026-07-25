---
name: changelog
description: >-
  Maintain CHANGELOG.md as a concise record of user-visible CLI, SDK, hook
  protocol, and security behavior. Use when adding release notes, documenting
  shipped user-facing features, or editing the [Unreleased] section.
---

# Changelog

Follow [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the policy
in [CONTRIBUTING.md](../../../CONTRIBUTING.md).

## Include

Add an entry when users or SDK consumers can observe a new or changed:

- CLI command, flag, discovery rule, installation effect, cache rule, report,
  or exit behavior;
- public SDK event, registration, builder, tool input, or compatibility
  contract;
- native protocol mapping or output behavior;
- security property.

Write outcomes, not implementation history:

```markdown
- `wat install` removes stale wat-managed events while preserving unrelated
  native hook entries.
```

## Exclude

Do not mention scaffolding of repository packages, internal refactors, test
harnesses, CI/lint changes, agent rules/skills, or documentation-only work.

Keep related details under one feature bullet. Avoid listing internal type
aliases, codec files, and helper packages unless they are themselves public API
that users call.

## Sections

Use `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, and `Security`.
Before a published baseline, describe initial-release behavior under `Added`;
do not claim a behavior changed or was removed from users when it was never
released.

If a change has no user-visible behavior, leave the changelog untouched.
