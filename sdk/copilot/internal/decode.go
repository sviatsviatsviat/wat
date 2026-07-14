package internal

// Decoder parses raw JSON into a decoded event value.
type Decoder func(raw []byte, received, canonical string) (any, error)

var decoders = map[string]Decoder{}

// RegisterDecoder registers a hook event decoder by canonical event name.
func RegisterDecoder(name string, fn Decoder) {
	decoders[name] = fn
}

// DecoderFor returns the decoder registered for name.
func DecoderFor(name string) (Decoder, bool) {
	fn, ok := decoders[name]
	return fn, ok
}

// RegisteredDecoders returns the canonical event names with registered decoders.
func RegisteredDecoders() []string {
	out := make([]string, 0, len(decoders))
	for name := range decoders {
		out = append(out, name)
	}
	return out
}
