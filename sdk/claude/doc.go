// Package claude is the Claude Code hook SDK. Hook authors register typed
// handlers with Chain methods into the shared run registry, then call run.Main
// from github.com/sviatsviatsviat/wat/sdk/run.
//
// Handlers build responses with hook-scoped *Results builders (and With* for
// advanced fields). A nil return is a silent no-op. Encode returns JSON
// only (no process exit code). Blocking is expressed via output fields; handler
// errors exit 1 by default or 2 with WithFailPolicy(FailBlock).
//
// See ExampleChain in example_test.go for a minimal handler.
package claude
