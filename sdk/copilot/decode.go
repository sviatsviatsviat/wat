package copilot

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

var codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)

type decodeFn func([]byte) (Event, error)

func registerDecoder(name string, fn decodeFn) {
	codec.Register(name, func(raw []byte) (any, error) {
		return fn(raw)
	})
}

func decodeAs[T Event](raw []byte) (Event, error) {
	return decodeAsAndThen[T](raw, nil)
}

func decodeAsAndThen[T Event](raw []byte, after func(*T, []byte)) (Event, error) {
	ev, err := hookkit.DecodeAsAndThen(raw, after)
	if err != nil {
		return nil, codec.WrapDecodeError(ev, err)
	}
	return ev, nil
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return ev.envelope()
}
