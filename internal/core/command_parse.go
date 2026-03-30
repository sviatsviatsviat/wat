package core

// Shell dialect labels returned by parsers and dialect detection (e.g. in internal/shellparse).
const (
	// DialectBash is the dialect label for POSIX/Bash-style parsing.
	DialectBash = "bash"
	// DialectPowerShell is the dialect label for PowerShell parsing.
	DialectPowerShell = "powershell"
)

// CommandNode is one simple shell command stage (after pipeline splitting where applicable).
// It is host-neutral: produced by bash or PowerShell parsers for guards and other consumers.
type CommandNode struct {
	// Name is the invoked command (e.g. executable or cmdlet name), after alias resolution on the PowerShell fast path.
	Name string
	// Args are positional arguments only (not flag values).
	Args []string
	// Flags maps parameter names to values for named parameters (e.g. "-Path" -> "C:\\temp").
	Flags map[string]string
	// Switches lists boolean or valueless parameters (e.g. "-Force", "-Recurse").
	Switches []string
	// PipeIndex is the zero-based index of this stage within its pipeline.
	PipeIndex int
	// PipeLength is the number of stages in that pipeline.
	PipeLength int
}

// ParseResult is the outcome of parsing a full shell command string.
type ParseResult struct {
	// Dialect is [DialectBash] or [DialectPowerShell], matching which parser produced the result.
	Dialect string
	// Pipeline is the ordered list of simple commands extracted from the input.
	Pipeline []CommandNode
	// Raw is the original command string passed to [ShellParser.Parse].
	Raw string
}

// ShellParser parses shell input for one dialect (see internal/shellparse).
type ShellParser interface {
	// Dialect returns [DialectBash] or [DialectPowerShell].
	Dialect() string
	// Parse analyzes raw and returns a [ParseResult] or a parse error (dialect-dependent).
	Parse(raw string) (ParseResult, error)
}
