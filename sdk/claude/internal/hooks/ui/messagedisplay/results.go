package messagedisplay

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Override returns a display-content override result.
	Override(content string) Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Override returns a display-content override result.
func (results) Override(content string) Output {
	c := content
	return output{displayContent: &c}
}
