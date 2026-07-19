// Package cursor is the Cursor hook SDK. Hook authors register typed handlers
// with On* helpers into the shared run registry, then call run.Main from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// Permission-gating events return exit code 2 on deny. Handler errors exit 1
// under Cursor's default fail-open policy. Hook stdout JSON is encoded
// internally from Output (sealed; only this package implements it). Payloads
// must include hook_event_name on the wire.
package cursor
