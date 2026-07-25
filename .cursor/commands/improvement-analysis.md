# Improvement analysis

Review the current conversation (or the described session): corrections, fixes, struggles, and review loops. Produce actionable process improvements for this repository.

## Steps

1. **Inventory** — List what was implemented, what was fixed in follow-up passes, and what failed or needed rework (tests, CI, docs, ambiguous specs).
2. **Root cause** — For each struggle, identify whether the gap was: missing rule/skill, underspecified task, environment (e.g. missing devcontainer), or review timing.
3. **Recommend** — Propose improvements ordered by leverage. Prefer updating `.cursor/skills/`, `.cursor/rules/`, `AGENTS.md`, `docs/`, or task DoD templates over one-off fixes.
4. **Triage** — Mark each recommendation **apply now**, **defer**, or **skip** with a one-line reason.
5. **Do not implement** unless the user asks — this command is analysis only.

## Output format

### What went well
- …

### Struggles and root causes
| Issue | Root cause |
|-------|------------|
| … | … |

### Recommendations (priority order)
1. **Title** — What to change, where (file/skill), and why it prevents recurrence.

### Suggested next actions
- Bullet list the user can approve for Agent mode.

## Reference

- Architecture and extension patterns: `docs/architecture.md`
- Public API: `docs/sdk.md`
- CLI behavior: `docs/usage.md`
- Contribution policy: `CONTRIBUTING.md`
- Task planning DoD: `.cursor/skills/task-planning/SKILL.md`
- Tests: `.cursor/skills/go-tests/SKILL.md`
- CI pins: `.cursor/skills/ci-pins/SKILL.md`
