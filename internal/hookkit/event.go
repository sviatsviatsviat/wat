package hookkit

// Event is implemented by every decoded per-agent hook event.
type Event interface {
	EventName() string
}

// Decoder parses raw JSON into a decoded event value.
type Decoder func(raw []byte) (Event, error)

// Output is a hook response. Concrete per-agent types implement IsZero, Encode,
// Merge, and Stop.
type Output interface {
	IsZero() bool
	Encode() (stdout []byte, exit int, err error)
	// Merge combines other into the receiver. other must be the same concrete type.
	// warnings lists non-fatal issues (e.g. last-wins overwrite of updatedInput).
	Merge(other Output) (merged Output, warnings []string, err error)
	// Stop reports whether remaining handlers should be skipped.
	Stop() bool
}

// BodyOnStderr is an optional Output extension. When BodyOnStderr returns true,
// run writes Encode's body to stderr instead of stdout. Native outputs use this
// when a host ignores stdout for a given exit code and reads the reason from
// stderr (for example Claude BlockExit).
type BodyOnStderr interface {
	BodyOnStderr() bool
}
