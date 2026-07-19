package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is any Cursor hook response. Only this package implements it.
type Output interface {
	isCursorOutput()
}

// outputEncoder is implemented by concrete hook outputs for encode.
type outputEncoder interface {
	Output
	isZero() bool
	allowedEvents() []string
	encode(eventName string) ([]byte, int, error)
}

// encode renders a typed output as Cursor stdout JSON and returns the
// process exit code.
func encode(eventName string, out Output) ([]byte, int, error) {
	normalized := hookkit.NormalizeOutput(out)
	if normalized == nil {
		return nil, 0, nil
	}
	enc, ok := normalized.(outputEncoder)
	if !ok {
		return nil, 0, fmt.Errorf("cursor: encode: unsupported output type %T", normalized)
	}
	if enc.isZero() {
		return nil, 0, nil
	}
	if err := hookkit.ValidateEncodePair(Dialect, eventName, normalized, enc.allowedEvents(), nil); err != nil {
		return nil, 0, err
	}
	return enc.encode(eventName)
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
