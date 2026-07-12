# Agent hook formats

Reference for tool names, MCP naming, and payload conventions across Claude Code, GitHub Copilot, and Cursor. Use this when implementing normalization, codecs, matchers, or portconfig translation.

Sources: [GitHub Copilot hooks reference](https://docs.github.com/en/copilot/reference/hooks-reference), [copilot-sdk#869](https://github.com/github/copilot-sdk/issues/869).

## Tool name normalization

`agenthooks.NormalizeToolName` maps native names to a canonical vocabulary. `ToolCall.Native` always keeps the original string.

| Agent | Surface | Builtin example | Normalized | MCP example | MCP detection |
|-------|---------|-----------------|------------|-------------|---------------|
| Claude Code | `PreToolUse` | `Bash` | `bash` | `mcp__github__create_issue` | `mcp__` prefix → `mcp=true`, name unchanged |
| Copilot | PascalCase `PreToolUse` | `Bash` | `bash` | `mcp__github__create_issue` | same as Claude (Claude-format matchers) |
| Copilot | camelCase `preToolUse` | `bash` | `bash` | `my-server-list_items` | codec sets `ToolCall.MCP` from structured metadata (not inferred from hyphens) |
| Cursor | `preToolUse` / dedicated shell hooks | `Shell` | `bash` | `MCP:browser_navigate` | `MCP:` prefix → `mcp=true`, name unchanged |

### Builtin alias examples

| Native (any case where noted) | Normalized |
|-------------------------------|------------|
| `Bash`, `bash`, `Shell`, `powershell` | `bash` |
| `Edit`, `edit`, `notebookedit` | `edit` |
| `Write`, `create` | `write` |
| `Read`, `view` | `read` |
| `Agent`, `task` | `task` |
| `web_fetch` | `web_fetch` |

Unknown names pass through unchanged with `mcp=false` unless an `mcp__` or `MCP:` prefix matches.

### Out of scope for `NormalizeToolName`

- **Copilot `preMcpToolCall`** — separate hook event with `serverName` and bare `toolName` fields, not a combined tool name string.
- **Cursor dedicated events** (`beforeShellExecution`, `beforeMCPExecution`, …) — folded into unified `Kind` by codecs; native event name stays in `Event.Name`.

## Unified event envelope

All codecs produce `agenthooks.Event`:

- `Agent` — `Dialect` (Claude, Copilot, Cursor)
- `Kind` — normalized category (`KindPreTool`, `KindStop`, …)
- `Name` — native hook event name as received
- `Raw` — untouched native JSON payload

See `go doc github.com/sviatsviatsviat/wat/agenthooks Event` for sub-structs (`Tool`, `Result`, `Life`, …).

## Related code

- Normalization: [`agenthooks/event.go`](../agenthooks/event.go) — `NormalizeToolName`, `InputAs`
- Tests: [`agenthooks/event_test.go`](../agenthooks/event_test.go)
