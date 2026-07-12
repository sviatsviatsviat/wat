# wat documentation

User-facing reference for the redesigned **wat** monorepo. Update these files when shipped behavior changes — not only CHANGELOG.md.

## Index

| Document | Contents |
|----------|----------|
| [agent-formats.md](agent-formats.md) | Tool names, MCP naming, and normalization rules per agent |
| [../README.md](../README.md) | Build, test, and current shipped capabilities |
| [../CHANGELOG.md](../CHANGELOG.md) | Release notes ([Unreleased] during development) |
| [../AGENTS.md](../AGENTS.md) | Agent/CI conventions for contributors |

## When to update

| Change | Update |
|--------|--------|
| New public `agenthooks` API | CHANGELOG, godoc, agent-formats if tool/dialect behavior; README only if build/docs links change |
| New `wat` subcommand | README, CHANGELOG, CLI tests per go-tests skill |
| New agent codec or MCP rule | agent-formats.md, normalization tests |
| CI action bump | AGENTS.md / README if versions cited; ci-pins skill audit |

Design tasks and prototypes stay in local-only `plan/` (gitignored).
