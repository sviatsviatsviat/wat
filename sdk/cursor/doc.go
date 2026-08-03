// Package cursor is the Cursor hook SDK. Hook authors register typed handlers
// with UseHooks, then pass the chain to run.Serve from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// Permission-gating events typically return exit code 2 on deny. SubagentStart,
// BeforeReadFile, and BeforeTabFileRead deny use exit 0 with permission JSON so
// Cursor applies the schema field. BeforeReadFile uses user_message (no ask or
// agent_message). BeforeTabFileRead encodes allow|deny only (no ask or message
// fields). BeforeShellExecution and BeforeMCPExecution enforce permission ask
// (user approval); PreToolUseResults.Ask encodes ask but Cursor does not enforce
// it today (prefer Allow or Deny); SubagentStart and BeforeReadFile expose
// SoftDeny (same encoding as Deny) instead of Ask. BeforeSubmitPrompt blocks
// with continue:false and exit 0 (not exit 2); its hooks.json matcher value is
// UserPromptSubmit. Handler errors exit 1 under Cursor's default fail-open
// policy; for security-critical BeforeMCPExecution gates, set failClosed: true
// in hooks.json so crash/timeout/invalid JSON blocks the tool.
// BeforeMCPExecution is deferred / not available for cloud agents. Hook stdout
// JSON is encoded internally from Output (sealed; only this package implements
// it). Payloads must include hook_event_name on the wire.
//
// AfterShellExecution, AfterMCPExecution, AfterFileEdit, and PostToolUseFailure
// are observe-only (Cursor Hooks docs list no consumed output fields). Rewrite
// MCP tool output with PostToolUse. Cloud agents do not load MCP hooks.
//
// Live subagentStart payloads may use automatic model sentinels
// ("", "auto", "default", "inherit") and kebab-case subagent_type values;
// see field godoc on SubagentStart and docs/agent-formats.md.
//
// Tool input on PreToolUse and related events is typed as [Input] with AsShell,
// AsRead, and related accessors (and Tool* name constants) in this package.
//
// See ExampleUseHooks in example_test.go for a runnable example.
package cursor
