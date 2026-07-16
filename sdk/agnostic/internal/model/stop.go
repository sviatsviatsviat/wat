package model

// StopResult is the portable hook response for Stop and SubagentStop events.
// Construct via StopFollowUp or agnostic.StopResults.
// A nil value is a no-op.
type StopResult interface {
	isStopResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// Result converts into the wire Result type for codecs.
	Result() Result
}

type stopResult struct {
	followUp string
}

func (stopResult) isStopResult() {}

// IsZero reports whether the result carries no instruction.
func (r stopResult) IsZero() bool { return r.followUp == "" }

// Result converts r into the wire Result type for codecs.
func (r stopResult) Result() Result { return Result{FollowUp: r.followUp} }

// StopFollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
func StopFollowUp(text string) StopResult { return stopResult{followUp: text} }
