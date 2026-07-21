package hookkit

import "fmt"

// Event is implemented by every decoded per-agent hook event.
type Event interface {
	EventName() string
}

// Encoder renders a typed hook output as native stdout JSON and an exit code.
type Encoder interface {
	Encode(eventName string, out Output) (stdout []byte, exit int, err error)
}

// Decoder parses raw JSON into a decoded event value.
type Decoder func(raw []byte) (Event, error)

// Codec owns a per-agent decoder registry, optional encoder, and dialect-scoped error sentinels.
// Each agent SDK should hold its own Codec instance; registries must not be shared
// across dialects (wire names overlap with different Go types).
type Codec struct {
	dialect                string
	empty, decode, nameReq error
	m                      map[string]Decoder
	encoder                Encoder
}

// NewCodec returns an empty codec for dialect (for example "claude").
// empty, decode, and nameReq are the caller's sentinel errors.
// enc may be nil when the codec is decode-only (tests); Encode then returns an error.
func NewCodec(dialect string, empty, decode, nameReq error, enc Encoder) *Codec {
	return &Codec{
		dialect: dialect,
		empty:   empty,
		decode:  decode,
		nameReq: nameReq,
		m:       make(map[string]Decoder),
		encoder: enc,
	}
}

// Encode renders out for eventName using the codec's encoder.
func (c *Codec) Encode(eventName string, out Output) ([]byte, int, error) {
	if c.encoder == nil {
		return nil, 0, fmt.Errorf("%s: encode: no encoder", c.dialect)
	}
	return c.encoder.Encode(eventName, out)
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
