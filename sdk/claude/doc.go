// Package claude is the Claude Code hook SDK. Hook authors register typed
// handlers with UseHooks, then pass the chain to run.Serve from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// Handlers build responses with hook-scoped *Results builders (and With* for
// advanced fields). A nil return is a silent no-op. Hook stdout JSON is encoded
// internally from Output (sealed; only this package implements it). Blocking is
// expressed via output fields; handler errors exit 1.
//
// Tool input on PreToolUse and related events is typed as [Input] with AsBash,
// AsWrite, and related accessors (and Tool* name constants) in this package.
//
// See ExampleUseHooks in example_test.go for a minimal handler.
package claude
