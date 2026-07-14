package portconfig

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/claude"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/copilot"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
)

// Config is a normalized hook configuration independent of any single agent's
// config file shape.
type Config = model.Config

// Entry is one normalized hook registration.
type Entry = model.Entry

// NativeEntry preserves an unmappable native hook block for round-trip emit.
type NativeEntry = model.NativeEntry

// Warning describes information lost or approximated during parse or emit.
type Warning = model.Warning

// Parse reads a native hook configuration file for dialect and returns a
// normalized Config.
func Parse(data []byte, dialect agnostic.Dialect) (Config, []Warning, error) {
	switch dialect {
	case agnostic.Claude:
		return claude.Parse(data)
	case agnostic.Copilot:
		return copilot.Parse(data)
	case agnostic.Cursor:
		return cursor.Parse(data)
	default:
		return Config{}, nil, fmt.Errorf("portconfig: unsupported dialect %q", dialect)
	}
}

// Emit renders cfg as a native hook configuration file for dialect.
func Emit(cfg Config, dialect agnostic.Dialect) ([]byte, []Warning, error) {
	switch dialect {
	case agnostic.Claude:
		return claude.Emit(cfg)
	case agnostic.Copilot:
		return copilot.Emit(cfg)
	case agnostic.Cursor:
		return cursor.Emit(cfg)
	default:
		return nil, nil, fmt.Errorf("portconfig: unsupported dialect %q", dialect)
	}
}
