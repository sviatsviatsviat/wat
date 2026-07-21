package hookkit

import "fmt"

// Output is a hook response. Concrete per-agent types implement IsZero,
// AllowedEvents, and Encode; package encoders orchestrate validation and
// dialect-specific side effects before calling Encode.
type Output interface {
	IsZero() bool
	AllowedEvents() []string
	Encode(eventName string) (stdout []byte, exit int, err error)
}

// ValidateEncodePair checks that eventName is allowed for the given output.
// canonicalize may remap eventName before comparison; pass nil to compare directly.
func ValidateEncodePair(label, eventName string, out Output, canonicalize func(string) (string, bool)) error {
	if eventName == "" {
		return nil
	}
	canonical := eventName
	if canonicalize != nil {
		if mapped, ok := canonicalize(eventName); ok {
			canonical = mapped
		}
	}
	for _, name := range out.AllowedEvents() {
		if canonical == name {
			return nil
		}
	}
	return fmt.Errorf("%s: encode: event %q incompatible with output type %T", label, eventName, out)
}
