package portconfig

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
)

func warnf(format string, args ...any) Warning { return Warning(fmt.Sprintf(format, args...)) }

// Parse reads a native hook configuration file for dialect and returns a
// normalized Config.
func Parse(data []byte, dialect agnostic.Dialect) (Config, []Warning, error) {
	switch dialect {
	case agnostic.Claude:
		return parseClaude(data)
	case agnostic.Copilot:
		return parseCopilot(data)
	case agnostic.Cursor:
		return parseCursor(data)
	default:
		return Config{}, nil, fmt.Errorf("portconfig: unsupported dialect %q", dialect)
	}
}

// Emit renders cfg as a native hook configuration file for dialect.
func Emit(cfg Config, dialect agnostic.Dialect) ([]byte, []Warning, error) {
	switch dialect {
	case agnostic.Claude:
		return emitClaude(cfg)
	case agnostic.Copilot:
		return emitCopilot(cfg)
	case agnostic.Cursor:
		return emitCursor(cfg)
	default:
		return nil, nil, fmt.Errorf("portconfig: unsupported dialect %q", dialect)
	}
}
