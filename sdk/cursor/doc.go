// Package cursor is the Cursor hook SDK. Hook authors register typed handlers
// with UseHooks, then pass the chain to run.Serve from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// Permission-gating events typically return exit code 2 on deny. SubagentStart
// deny uses exit 0 with permission JSON so Cursor applies the schema field.
// BeforeShellExecution and BeforeMCPExecution enforce permission ask (user
// approval); PreToolUse accepts ask in the schema but does not enforce it, and
// SubagentStart treats ask as deny. Handler errors exit 1 under Cursor's default
// fail-open policy. Hook stdout JSON is encoded internally from Output (sealed;
// only this package implements it). Payloads must include hook_event_name on
// the wire.
//
// Tool input on PreToolUse and related events is typed as [Input] with AsShell,
// AsRead, and related accessors (and Tool* name constants) in this package.
//
// See ExampleUseHooks in example_test.go for a runnable example.
package cursor
