package cli

// RootHelpSummary is the canonical root usage text (commands and flags) for the wat binary.
const RootHelpSummary = `Wat runs templated hook subprocesses for agent hosts.

It reads hook event JSON from stdin, resolves argv placeholders, runs the
child process, and writes the hook protocol response to stdout.

Usage:

	wat <command> [templated arguments]

Supported commands:

	run        Run a templated hook subprocess

Flags:

	-H, --host <name>    Hook host that handles stdin and hook protocol
	                      output (default: cursor)

If equivalent flags repeated, the last value wins.`

const (
	rootHelpText = RootHelpSummary + `

Run wat with no arguments to print this text; run wat run without a subprocess command to print run usage.`

	runHelpText = `Usage:

	wat run <command> [templated arguments]
	wat run --host <name> <command> [templated arguments]
	wat run -H <name> <command> [templated arguments]

Host flags (after run, before the subprocess command) match root usage:

	-H, --host <name>    Hook host that handles stdin and hook protocol
	                      output (default: cursor)

If equivalent flags repeated, the last value wins.

The hook JSON on stdin supplies template values. Only these placeholders are allowed:

	__CONVERSATION_ID__  __GENERATION_ID__  __MODEL__
	__HOOK_EVENT_NAME__  __CURSOR_VERSION__  __WORKSPACE_ROOTS__
	__USER_EMAIL__       __TRANSCRIPT_PATH__

Wat prints {} on stdout for Cursor; the child's stderr is copied to wat's stderr
(child stdout is discarded — redirect with 1>&2 or >&2 if you need logs).

Examples:

Windows:

	wat run cmd /c "go version 1>&2"
	wat run cmd /c "echo __HOOK_EVENT_NAME__ 1>&2"
	wat run -H cursor cmd /c "go version 1>&2"

Unix / macOS:

	wat run sh -c "go version 1>&2"
	wat run sh -c 'echo __HOOK_EVENT_NAME__ 1>&2'
	wat run -H cursor sh -c "go version 1>&2"`
)
