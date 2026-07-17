package copilot

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// outputEncoder is implemented by concrete hook outputs for Encode.
type outputEncoder interface {
	isZero() bool
	allowedEvents() []string
	encode() ([]byte, int, error)
}

// Encode renders a typed output struct as Copilot flat camelCase stdout JSON and
// returns the process exit code.
func Encode(eventName string, out any) ([]byte, int, error) {
	out = hookkit.NormalizeOutput(out)
	if out == nil {
		return nil, 0, nil
	}
	enc, ok := out.(outputEncoder)
	if !ok {
		return nil, 0, fmt.Errorf("copilot: encode: unsupported output type %T", out)
	}
	if enc.isZero() {
		return nil, 0, nil
	}
	if err := hookkit.ValidateEncodePair(Dialect, eventName, out, enc.allowedEvents(), func(name string) (string, bool) {
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

func isZeroOutput(out any) bool {
	if z, ok := out.(interface{ isZero() bool }); ok {
		return z.isZero()
	}
	return hookkit.IsZeroOutput(out)
}

// IsZeroOutput reports whether out is an empty hook response.
func IsZeroOutput(out any) bool { return isZeroOutput(out) }
