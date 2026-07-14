package agnostic

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
)

// CodecFor returns the codec for a dialect.
func CodecFor(d Dialect) (Codec, error) {
	switch d {
	case Claude:
		return &claude.Codec{}, nil
	case Copilot:
		return &copilot.Codec{}, nil
	case Cursor:
		return &cursor.Codec{}, nil
	default:
		return nil, fmt.Errorf("agnostic: no codec for dialect %q", d)
	}
}
