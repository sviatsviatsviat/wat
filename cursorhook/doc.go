// Package cursorhook is the Cursor hook SDK. Hook authors register typed
// handlers; the SDK decodes stdin JSON, dispatches by event, encodes stdout,
// and selects exit codes per the Cursor hook protocol.
//
// Permission-gating events return exit code 2 on deny. Handler errors exit 1
// under Cursor's default fail-open policy. Payloads without hook_event_name
// require WithEvent when calling Decode or Serve.
package cursorhook
