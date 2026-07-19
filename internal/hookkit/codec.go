package hookkit

import "fmt"

// Event is implemented by every decoded per-agent hook event.
type Event interface {
	EventName() string
}

// Decoder parses raw JSON into a decoded event value.
type Decoder func(raw []byte) (Event, error)

// Codec owns a per-agent decoder registry and dialect-scoped error sentinels.
// Each agent SDK should hold its own Codec instance; registries must not be shared
// across dialects (wire names overlap with different Go types).
type Codec struct {
	dialect                string
	empty, decode, nameReq error
	m                      map[string]Decoder
}

// NewCodec returns an empty codec for dialect (for example "claude").
// empty, decode, and nameReq are the caller's sentinel errors.
func NewCodec(dialect string, empty, decode, nameReq error) *Codec {
	return &Codec{
		dialect: dialect,
		empty:   empty,
		decode:  decode,
		nameReq: nameReq,
		m:       make(map[string]Decoder),
	}
}

// Register associates wire event name with fn.
func (c *Codec) Register(name string, fn Decoder) {
	c.m[name] = fn
}

// EventName peeks hook_event_name from raw.
func (c *Codec) EventName(raw []byte) (string, error) {
	return RequireHookEventName(raw, c.empty, c.decode, c.nameReq)
}

// Decode peeks hook_event_name, looks up a registered decoder, and runs it.
// It returns an error if name is non-empty but not registered (Serve must not
// decode when no handlers exist for the event).
func (c *Codec) Decode(raw []byte) (Event, error) {
	name, err := c.EventName(raw)
	if err != nil {
		return nil, err
	}
	fn, ok := c.m[name]
	if !ok {
		return nil, fmt.Errorf("%s: decode: unknown hook event %s", c.dialect, name)
	}
	return fn(raw)
}

// WrapDecodeError wraps a JSON unmarshal failure with dialect and ErrDecodePayload.
func (c *Codec) WrapDecodeError(ev any, err error) error {
	return fmt.Errorf("%s: decode %T: %w", c.dialect, ev, fmt.Errorf("%w: %w", c.decode, err))
}
