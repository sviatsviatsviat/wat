package cli

// RootHelpSummary is the canonical root usage text (commands and flags) for the wat binary.
const RootHelpSummary = `Wat runs templated hook subprocesses for agent hosts.

It reads hook event JSON from stdin, resolves placeholders in the command template, runs the
child process, and writes the hook protocol response to stdout.

Usage:

	wat <host> <command> [arguments…]

The first word is the hook host (e.g. cursor). The second is the wat subcommand.

Supported commands:

	run        Run a templated hook subprocess

For run, optional flags before the subprocess template:

	-f, --file-pattern <re>     Optional; when stdin supplies __FILE_PATH__ (Cursor
	                             afterFileEdit), skip the subprocess if the path does
	                             not match <re> (Go regexp). Default * means no filter.
	                             If set, <re> must be non-empty.

If equivalent flags repeat, the last value wins.`

const (
	rootHelpText = RootHelpSummary + `

Run wat with no arguments to print this text; run wat <host> run without a subprocess command to print run usage.`

	runHelpText = `Usage:

	wat <host> run <command> [templated arguments]
	wat <host> run [-f <re>] <command> [templated arguments]
	wat <host> run [--file-pattern <re>] <command> [templated arguments]
	wat <host> run [--file-pattern=<re>] <command> [templated arguments]

Put -f/--file-pattern (if any) after run and before the subprocess command. If equivalent flags repeat, the last value wins.

When -f/--file-pattern is not the default (*), and the hook bindings include __FILE_PATH__, the subprocess runs only if the cleaned path matches the regexp.

The hook JSON on stdin supplies template values. Only these placeholders are allowed:

	__CONVERSATION_ID__  __GENERATION_ID__  __MODEL__
	__HOOK_EVENT_NAME__  __CURSOR_VERSION__  __WORKSPACE_ROOTS__
	__USER_EMAIL__       __TRANSCRIPT_PATH__  __FILE_PATH__

Wat prints {} on stdout for Cursor; the child's stderr is copied to wat's stderr
(child stdout is discarded — redirect with 1>&2 or >&2 if you need logs).

Examples:

Windows:

	wat cursor run cmd /c "go version 1>&2"
	wat cursor run cmd /c "echo __HOOK_EVENT_NAME__ 1>&2"

Unix / macOS:

	wat cursor run sh -c "go version 1>&2"
	wat cursor run sh -c 'echo __HOOK_EVENT_NAME__ >&2'`
)
