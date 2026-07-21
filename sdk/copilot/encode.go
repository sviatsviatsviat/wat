package copilot

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is any Copilot hook response.
type Output = hookkit.Output

type encoder struct{}

func newEncoder() hookkit.Encoder {
	return encoder{}
}

// Encode validates out and renders Copilot snake_case stdout JSON.
func (encoder) Encode(eventName string, out Output) ([]byte, int, error) {
	if out.IsZero() {
		return nil, 0, nil
	}
	if err := hookkit.ValidateEncodePair(Dialect, eventName, out, func(name string) (string, bool) {
		canonical, known := CanonicalEventName(name)
		if !known {
			return name, true
		}
		return canonical, true
	}); err != nil {
		return nil, 0, err
	}
	return out.Encode(eventName)
}
