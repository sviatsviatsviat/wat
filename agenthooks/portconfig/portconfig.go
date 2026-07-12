package portconfig

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/agenthooks"
)

func warnf(format string, args ...any) Warning { return Warning(fmt.Sprintf(format, args...)) }

// Parse reads a native hook configuration file for dialect and returns a
// normalized Config.
func Parse(data []byte, dialect agenthooks.Dialect) (Config, []Warning, error) {
	switch dialect {
	case agenthooks.Claude:
		return parseClaude(data)
	case agenthooks.Copilot:
		return parseCopilot(data)
	case agenthooks.Cursor:
		return parseCursor(data)
	default:
		return Config{}, nil, fmt.Errorf("portconfig: unsupported dialect %q", dialect)
	}
}

// Emit renders cfg as a native hook configuration file for dialect.
func Emit(cfg Config, dialect agenthooks.Dialect) ([]byte, []Warning, error) {
	switch dialect {
	case agenthooks.Claude:
		return emitClaude(cfg)
	case agenthooks.Copilot:
		return emitCopilot(cfg)
	case agenthooks.Cursor:
		return emitCursor(cfg)
	default:
		return nil, nil, fmt.Errorf("portconfig: unsupported dialect %q", dialect)
	}
}
