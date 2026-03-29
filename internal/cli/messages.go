package cli

// RootHelpSummary is the canonical root usage text (commands and flags) for the wat binary.
const RootHelpSummary = `Wat runs templated hook subprocesses for agent hosts.

It reads hook event JSON from stdin, resolves placeholders in the command template, runs the
child process, and writes the hook protocol response to stdout.

Usage:

	wat <host> <command> [arguments…]

The first word is the hook host (e.g. cursor). The second is the wat subcommand.

Supported commands:

	exec       Run a templated hook subprocess

For exec, optional flags before the subprocess template:

	-f, --file-pattern <re>     Optional; when stdin supplies __FILE_PATH__ (Cursor
	                             afterFileEdit or afterTabFileEdit), skip the subprocess if the path does
	                             not match <re> (Go regexp). Default * means no filter.
	                             If set, <re> must be non-empty.

If equivalent flags repeat, the last value wins.`

const (
	rootHelpText = RootHelpSummary + `

Run wat with no arguments to print this text; run wat <host> exec without a subprocess command to print exec usage.`

	execHelpText = `Usage:

	wat <host> exec <command> [templated arguments]
	wat <host> exec [-f <re>] <command> [templated arguments]
	wat <host> exec [--file-pattern <re>] <command> [templated arguments]
	wat <host> exec [--file-pattern=<re>] <command> [templated arguments]

Put -f/--file-pattern (if any) after exec and before the subprocess command. If equivalent flags repeat, the last value wins.

When -f/--file-pattern is not the default (*), and the hook bindings include __FILE_PATH__, the subprocess runs only if the cleaned path matches the regexp.

The hook JSON on stdin supplies template values. Only these placeholders are allowed:

	__CONVERSATION_ID__  __GENERATION_ID__  __MODEL__
	__HOOK_EVENT_NAME__  __CURSOR_VERSION__  __USER_EMAIL__
	__TRANSCRIPT_PATH__  __FILE_PATH__  __TOOL_NAME__  __DURATION__  __SANDBOX__  __DURATION_MS__
	__SESSION_ID__  __REASON__  __IS_BACKGROUND__  __FINAL_STATUS__  __ERROR_MESSAGE__

Wat prints {} on stdout for Cursor; the child's stderr is copied to wat's stderr
(child stdout is discarded — redirect with 1>&2 or >&2 if you need logs).

Examples:

Windows:

	wat cursor exec cmd /c "go version 1>&2"
	wat cursor exec cmd /c "echo __HOOK_EVENT_NAME__ 1>&2"

Unix / macOS:

	wat cursor exec sh -c "go version 1>&2"
	wat cursor exec sh -c 'echo __HOOK_EVENT_NAME__ >&2'`
)
