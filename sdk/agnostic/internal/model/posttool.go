package model

// PostToolResult is the portable hook response for PostTool events.
// Construct via PostToolContext or agnostic.PostToolResults, then With*.
// A nil value is a no-op.
type PostToolResult interface {
	isPostToolResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// Result converts into the wire Result type for codecs.
	Result() Result
	// WithUpdatedOutput replaces tool result text when set.
	WithUpdatedOutput(output string) PostToolResult
}

type postToolResult struct {
	updatedOutput *string
	context       string
}

func (postToolResult) isPostToolResult() {}

// IsZero reports whether the result carries no instruction.
func (r postToolResult) IsZero() bool {
	return r.updatedOutput == nil && r.context == ""
}

// Result converts r into the wire Result type for codecs.
func (r postToolResult) Result() Result {
	return Result{UpdatedOutput: r.updatedOutput, Context: r.context}
}

// WithUpdatedOutput replaces tool result text when set.
func (r postToolResult) WithUpdatedOutput(output string) PostToolResult {
	r.updatedOutput = &output
	return r
}

// PostToolContext returns a context-injection-only PostTool result.
func PostToolContext(text string) PostToolResult { return postToolResult{context: text} }
