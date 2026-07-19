package copilot

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is any Copilot hook response. Only this package implements it.
type Output interface {
	isCopilotOutput()
}

// outputEncoder is implemented by concrete hook outputs for encode.
type outputEncoder interface {
	Output
	isZero() bool
	allowedEvents() []string
	encode() ([]byte, int, error)
}

// encode renders a typed output struct as Copilot snake_case stdout JSON and
// returns the process exit code.
func encode(eventName string, out Output) ([]byte, int, error) {
	normalized := hookkit.NormalizeOutput(out)
	if normalized == nil {
		return nil, 0, nil
	}
	enc, ok := normalized.(outputEncoder)
	if !ok {
		return nil, 0, fmt.Errorf("copilot: encode: unsupported output type %T", normalized)
	}
	if enc.isZero() {
		return nil, 0, nil
	}
	if err := hookkit.ValidateEncodePair(Dialect, eventName, normalized, enc.allowedEvents(), func(name string) (string, bool) {
		canonical, known := CanonicalEventName(name)
		if !known {
			return name, true
		}
		return canonical, true
	}); err != nil {
		return nil, 0, err
	}
	return enc.encode()
}

func isZeroOutput(out Output) bool {
	if out == nil {
		return true
	}
	if z, ok := out.(interface{ isZero() bool }); ok {
		return z.isZero()
	}
	return hookkit.IsZeroOutput(out)
}

// IsZeroOutput reports whether out is an empty hook response.
func IsZeroOutput(out Output) bool { return isZeroOutput(out) }
