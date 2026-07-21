package cursor

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is any Cursor hook response.
type Output = hookkit.Output

type encoder struct{}

func newEncoder() hookkit.Encoder {
	return encoder{}
}

// Encode validates out and renders Cursor stdout JSON.
func (encoder) Encode(eventName string, out Output) ([]byte, int, error) {
	if out.IsZero() {
		return nil, 0, nil
	}
	if err := hookkit.ValidateEncodePair(Dialect, eventName, out, nil); err != nil {
		return nil, 0, err
	}
	return out.Encode(eventName)
}
