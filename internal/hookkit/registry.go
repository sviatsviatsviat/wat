package hookkit

// Decoder parses raw JSON into a decoded event value.
// received is the wire event name; canonical is the normalized name used for lookup.
type Decoder func(raw []byte, received, canonical string) (any, error)

// DecoderRegistry maps canonical event names to decode functions.
// Each agent SDK should hold its own registry instance.
type DecoderRegistry struct {
	m map[string]Decoder
}

// NewDecoderRegistry returns an empty decoder registry.
func NewDecoderRegistry() *DecoderRegistry {
	return &DecoderRegistry{m: make(map[string]Decoder)}
}

// Register associates name with fn.
func (r *DecoderRegistry) Register(name string, fn Decoder) {
	r.m[name] = fn
}

// Lookup returns the decoder registered for name.
func (r *DecoderRegistry) Lookup(name string) (Decoder, bool) {
	fn, ok := r.m[name]
	return fn, ok
}
