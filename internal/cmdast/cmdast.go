// Package cmdast defines a lean dialect-agnostic AST for hook command lines
// (pipelines, chains, fully modeled simple commands, opaque compound constructs)
// and dialect sniffing ([DetectDialect], [HostHintFromGOOS]). To parse a hook command
// string end-to-end, use subpackage hookline’s ParseCommandLine.
package cmdast

// SourceLang identifies the shell dialect.
type SourceLang string

const (
	LangBash       SourceLang = "bash"
	LangPosix      SourceLang = "posix"
	LangPowerShell SourceLang = "powershell"
)

// Span marks a range in the original input.
type Span struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position is a byte offset and 1-based line/column (when available).
type Position struct {
	Offset uint `json:"offset"`
	Line   uint `json:"line"`
	Col    uint `json:"col"`
}

// CommandLine is the parsed result of a hook command string.
type CommandLine struct {
	Lang         SourceLang   `json:"lang"`
	Raw          string       `json:"raw"`
	Statements   []*Statement `json:"statements"`
	ParseError   string       `json:"parse_error,omitempty"`
	ParsePartial bool         `json:"parse_partial,omitempty"`
}

// StmtKind is the statement discriminator.
type StmtKind string

const (
	StmtCommand  StmtKind = "Command"
	StmtPipeline StmtKind = "Pipeline"
	StmtChain    StmtKind = "Chain"
	StmtCompound StmtKind = "Compound"
)

// Statement is a discriminated union; exactly one of the pointer fields is non-nil.
type Statement struct {
	Kind     StmtKind  `json:"kind"`
	Span     Span      `json:"span"`
	Command  *Command  `json:"command,omitempty"`
	Pipeline *Pipeline `json:"pipeline,omitempty"`
	Chain    *Chain    `json:"chain,omitempty"`
	Compound *Compound `json:"compound,omitempty"`
}

// Command is a simple invocation (one stage of a pipeline or chain).
type Command struct {
	Name        string     `json:"name"`
	Args        []Arg      `json:"args,omitempty"`
	Redirects   []Redirect `json:"redirects,omitempty"`
	Assignments []Assign   `json:"assignments,omitempty"`
	Background  bool       `json:"background,omitempty"`
	Negated     bool       `json:"negated,omitempty"`
}

// Arg is one argument token with optional structure for guards.
type Arg struct {
	Span     Span     `json:"span"`
	Literal  string   `json:"literal"`
	Flag     *Flag    `json:"flag,omitempty"`
	Vars     []VarRef `json:"vars,omitempty"`
	CmdSubs  []CmdSub `json:"cmd_subs,omitempty"`
	Expanded bool     `json:"expanded,omitempty"`
}

// Flag is a detected named parameter (POSIX-style or from PS parser).
type Flag struct {
	Name   string `json:"name"`
	Dashes int    `json:"dashes"`
	Value  *Arg   `json:"value,omitempty"`
}

// VarRef is a variable reference embedded in an argument.
type VarRef struct {
	Span  Span   `json:"span"`
	Name  string `json:"name"`
	Scope string `json:"scope,omitempty"`
}

// CmdSub is a command substitution with recursively parsed content.
type CmdSub struct {
	Span     Span         `json:"span"`
	Commands []*Statement `json:"commands"`
}

// Assign is a variable assignment.
type Assign struct {
	Span     Span   `json:"span"`
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    Arg    `json:"value"`
}

// Redirect is an I/O redirection.
type Redirect struct {
	Span     Span   `json:"span"`
	Operator string `json:"operator"`
	Target   string `json:"target"`
}

// Pipeline is cmd | cmd | …
type Pipeline struct {
	Stages     []*Statement `json:"stages"`
	Background bool         `json:"background,omitempty"`
}

// Chain is cmd && cmd, cmd || cmd, or cmd ; cmd (binary recursive tree).
type Chain struct {
	Operator string     `json:"operator"`
	Left     *Statement `json:"left"`
	Right    *Statement `json:"right"`
}

// CompoundKind classifies an opaque compound construct.
type CompoundKind string

const (
	CompoundIf         CompoundKind = "If"
	CompoundFor        CompoundKind = "For"
	CompoundForEach    CompoundKind = "ForEach"
	CompoundWhile      CompoundKind = "While"
	CompoundUntil      CompoundKind = "Until"
	CompoundCase       CompoundKind = "Case"
	CompoundSwitch     CompoundKind = "Switch"
	CompoundTry        CompoundKind = "Try"
	CompoundFunction   CompoundKind = "Function"
	CompoundSubshell   CompoundKind = "Subshell"
	CompoundBlock      CompoundKind = "Block"
	CompoundTest       CompoundKind = "Test"
	CompoundArithmetic CompoundKind = "Arithmetic"
	// CompoundOther covers shell constructs not enumerated above (e.g. time, coproc) while keeping Raw.
	CompoundOther CompoundKind = "Other"
)

// Compound flags a complex construct without modeling its body.
type Compound struct {
	CompoundKind CompoundKind `json:"compound_kind"`
	Raw          string       `json:"raw"`
}
