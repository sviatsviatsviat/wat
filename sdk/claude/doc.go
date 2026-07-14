// Package claude is the Claude Code hook SDK. Hook authors register typed
// handlers with On (chainable via Chain) into the shared run registry, then
// call run.Main from github.com/sviatsviatsviat/wat/sdk/run.
//
// Encode returns JSON only (no process exit code). Blocking is expressed via
// output fields; handler errors exit 1 by default or 2 with WithFailPolicy(FailBlock).
//
// See ExampleOn in example_test.go for a minimal handler.
package claude
