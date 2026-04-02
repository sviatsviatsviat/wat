// Package hookline wires dialect detection to the bash and PowerShell converters under cmdast.
package hookline

import (
	"github.com/sviatsviatsviat/wat/internal/cmdast"
	"github.com/sviatsviatsviat/wat/internal/cmdast/psconv"
	"github.com/sviatsviatsviat/wat/internal/cmdast/shconv"
)

// ParseCommandLine parses raw into a [cmdast.CommandLine] using [cmdast.DetectDialect] with [cmdast.HostHintFromGOOS].
// Bash-style input uses mvdan.cc/sh ([shconv.Parse]); PowerShell uses pwsh or powershell and an embedded script ([psconv.Parse])
// with the command on stdin.
//
// The returned dialect string is [cmdast.DialectBash] or [cmdast.DialectPowerShell].
func ParseCommandLine(raw string) (*cmdast.CommandLine, string, error) {
	dialect := cmdast.DetectDialect(raw, cmdast.HostHintFromGOOS())
	switch dialect {
	case cmdast.DialectPowerShell:
		cl, err := psconv.Parse(raw)
		return cl, dialect, err
	default:
		cl, err := shconv.Parse(raw, cmdast.LangBash)
		return cl, dialect, err
	}
}
