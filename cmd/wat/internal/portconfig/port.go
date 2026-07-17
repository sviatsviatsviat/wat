package portconfig

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/claude"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/copilot"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
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
func Parse(data []byte, dialect string) (Config, []Warning, error) {
	switch dialect {
	case sdkclaude.Dialect:
		return claude.Parse(data)
	case sdkcopilot.Dialect:
		return copilot.Parse(data)
	case sdkcursor.Dialect:
		return cursor.Parse(data)
	default:
		return Config{}, nil, fmt.Errorf("portconfig: unsupported dialect %q", dialect)
	}
}

// Emit renders cfg as a native hook configuration file for dialect.
func Emit(cfg Config, dialect string) ([]byte, []Warning, error) {
	switch dialect {
	case sdkclaude.Dialect:
		return claude.Emit(cfg)
	case sdkcopilot.Dialect:
		return copilot.Emit(cfg)
	case sdkcursor.Dialect:
		return cursor.Emit(cfg)
	default:
		return nil, nil, fmt.Errorf("portconfig: unsupported dialect %q", dialect)
	}
}
