package shellparse

// PowerShell fast-path: [splitPipeline], [tokenizeStage], and [parseFast] delimiter / quote runes.
const (
	psRuneBacktick    = '`'
	psRuneSingleQuote = '\''
	psRuneDoubleQuote = '"'
	psRunePipe        = '|'
	psRuneSemicolon   = ';'
	psRuneSpace       = ' '
	psRuneTab         = '\t'
)

// PowerShell tokenizer emits these as separate string tokens.
const (
	psTokenPipe      = "|"
	psTokenSemicolon = ";"
)

// POSIX-style switch / long-option prefix for fast-path argument classification (bash and PowerShell).
const switchPrefix = "-"

// detectComplexity string probes — conservative signals to delegate to **pwsh** + embedded parser.
const (
	psComplexBraceOpen       = "{"
	psComplexSubExpr         = "$("
	psComplexArrayLiteral    = "@("
	psComplexHashtableLit    = "@{"
	psComplexHereStringDbl   = "@\""
	psComplexHereStringSgl   = "@'"
	psComplexParenRunes      = "()" // argument to [strings.ContainsAny]
	psComplexPipe            = "|"
	psComplexBacktickNewline = "`\n"
	psComplexNewline         = "\n"
)
