package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
)

// Unsupported lists Result capabilities that dialect d cannot express for
// event kind k. Codecs drop these fields silently; callers may warn.
func Unsupported(d Dialect, k Kind, r Result) []string {
	switch d {
	case Claude:
		return claude.Unsupported(k, r)
	case Copilot:
		return copilot.Unsupported(k, r)
	case Cursor:
		return cursor.Unsupported(k, r)
	default:
		return nil
	}
}
