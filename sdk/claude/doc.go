// Package claude is the Claude Code hook SDK. Hook authors register typed
// handlers with On* helpers into the shared run registry, then call run.Main
// from github.com/sviatsviatsviat/wat/sdk/run.
//
// Handlers build responses with hook-scoped *Results builders (and With* for
// advanced fields). A nil return is a silent no-op. Hook stdout JSON is encoded
// internally from Output (sealed; only this package implements it). Blocking is
// expressed via output fields; handler errors exit 1.
//
// See ExampleOnPreToolUse in example_test.go for a minimal handler.
package claude
