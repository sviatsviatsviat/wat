package model

// SessionStartResult is the portable hook response for SessionStart events.
// Construct via SessionStartContext or agnostic.SessionStartResults.
// A nil value is a no-op.
type SessionStartResult interface {
	isSessionStartResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// Result converts into the wire Result type for codecs.
	Result() Result
}

type sessionStartResult struct {
	context string
}

func (sessionStartResult) isSessionStartResult() {}

// IsZero reports whether the result carries no instruction.
func (r sessionStartResult) IsZero() bool { return r.context == "" }

// Result converts r into the wire Result type for codecs.
func (r sessionStartResult) Result() Result { return Result{Context: r.context} }

// SessionStartContext returns a context-injection-only SessionStart result.
func SessionStartContext(text string) SessionStartResult {
	return sessionStartResult{context: text}
}
