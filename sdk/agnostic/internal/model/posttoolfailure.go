package model

// PostToolFailureResult is the portable hook response for PostToolFailure events.
// Construct via PostToolFailureContext or agnostic.PostToolFailureResults.
// A nil value is a no-op.
type PostToolFailureResult interface {
	isPostToolFailureResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// Result converts into the wire Result type for codecs.
	Result() Result
}

type postToolFailureResult struct {
	context string
}

func (postToolFailureResult) isPostToolFailureResult() {}

// IsZero reports whether the result carries no instruction.
func (r postToolFailureResult) IsZero() bool { return r.context == "" }

// Result converts r into the wire Result type for codecs.
func (r postToolFailureResult) Result() Result { return Result{Context: r.context} }

// PostToolFailureContext returns recovery guidance for PostToolFailure events.
func PostToolFailureContext(text string) PostToolFailureResult {
	return postToolFailureResult{context: text}
}
