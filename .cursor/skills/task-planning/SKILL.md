---
name: task-planning
description: >-
  Plan wat implementation tasks with scoped work, dependencies, and a Definition
  of Done checklist. Use when creating or updating plan/tasks/*.md files or
  drafting implementation plans for design tasks.
---

# Task planning

Use when breaking down work from `plan/hooks-abstraction-design.md` into `plan/tasks/*.md` or when the user asks to plan a task.

## Task file structure

Each task file should include:

1. **Depends on** — upstream tasks that must land first
2. **Goal** — one paragraph outcome
3. **Work** — numbered implementation steps (scoped; defer out-of-scope items explicitly)
4. **Expected result** — acceptance criteria in prose
5. **Definition of done** — checkbox checklist agents run before marking complete

## Definition of done template

Copy and tailor this block into every new or updated task plan:

```markdown
## Definition of done

- [ ] Scope matches **Work** only; out-of-scope items listed explicitly
- [ ] Every new exported symbol has godoc (see godoc skill)
- [ ] Tests assert real behavior; table-driven where appropriate (see go-tests skill)
- [ ] `go vet ./...` and `go test ./...` pass for touched packages
- [ ] `golangci-lint run ./...` passes (includes gofmt)
- [ ] User-facing API/CLI changes have CHANGELOG **Added** entries (see changelog skill)
- [ ] User-facing behavior documented in [docs/](../../docs/) and [CHANGELOG.md](../../CHANGELOG.md) when shipped
- [ ] `.github/workflows/ci.yml` action pins are exact patch tags, not floating majors (see ci-pins skill)
- [ ] Agent/tool naming tasks reference [docs/agent-formats.md](../../docs/agent-formats.md)
```

### Tailor per task type

| Task type | Add to DoD |
|-----------|------------|
| `agenthooks` API | Normalization/codec tests per agent column in agent-formats doc |
| CLI (`cmd/wat`) | Help on stdout; invalid usage on stderr + exit 1; CLI table tests with `wantOutput` |
| Codec | Fixture decode test per supporting agent; `Event.Raw` preserved |
| CI / workflow | Every `uses:` line pinned; bump all references to same dependency together |
| Docs-only | No CHANGELOG entry unless user-facing behavior text changed |

## Planning workflow

1. Read the design section and any prototype under `plan/agenthooks/`.
2. List **in scope** vs **deferred** (name the future task file).
3. Ask the user when a cross-task dependency is ambiguous (e.g. minimal type stub now vs full feature later).
4. Add **Definition of done** before handing off to implementation.
5. Link to [docs/agent-formats.md](../../docs/agent-formats.md) when the task touches tool names, MCP, or dialect payloads.

## Anti-patterns

- Do not plan "move entire prototype" when the task only needs one package.
- Do not leave DoD implicit in **Expected result** — use an explicit checklist.
- Do not document MCP or tool naming from memory; verify against agent-formats doc or primary sources.
