// Package cursor is the Cursor hook SDK. Hook authors register typed handlers
// with UseHooks, then call run.Main from github.com/sviatsviatsviat/wat/sdk/run.
//
// Permission-gating events return exit code 2 on deny. Handler errors exit 1
// under Cursor's default fail-open policy. Hook stdout JSON is encoded
// internally from Output (sealed; only this package implements it). Payloads
// must include hook_event_name on the wire.
//
// Tool input on PreToolUse and related events is typed as [Input] with AsShell,
// AsRead, and related accessors (and Tool* name constants) in this package.
//
// See ExampleUseHooks in example_test.go for a runnable example.
package cursor
