package shellparse

import "github.com/sviatsviatsviat/wat/internal/core"

var (
	_ core.ShellParser = (*BashParser)(nil)
	_ core.ShellParser = (*PowerShellParser)(nil)
)

// ParserRouter selects bash vs PowerShell parsing from the command string.
type ParserRouter struct {
	bash       *BashParser
	powershell *PowerShellParser
}

// NewParserRouter returns a router with default parser implementations.
func NewParserRouter() *ParserRouter {
	return &ParserRouter{
		bash:       &BashParser{},
		powershell: &PowerShellParser{},
	}
}

// Parse detects dialect from markers and an OS-derived host hint when inconclusive, then runs the matching parser.
func (r *ParserRouter) Parse(raw string) (core.ParseResult, error) {
	dialect := detectDialect(raw, hostHintFromGOOS())
	switch dialect {
	case core.DialectPowerShell:
		return r.powershell.Parse(raw)
	default:
		return r.bash.Parse(raw)
	}
}
