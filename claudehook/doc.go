// Package claudehook is the Claude Code hook SDK. Hook authors register typed
// handlers with On and call Main; the SDK decodes stdin JSON, dispatches by
// event, encodes stdout, and selects exit codes per the Claude Code hook protocol.
//
// Encode returns JSON only (no process exit code). Blocking is expressed via
// output fields; handler errors exit 1 by default or 2 with WithFailPolicy(FailBlock).
//
// See ExampleMux in example_test.go for a minimal handler.
package claudehook
