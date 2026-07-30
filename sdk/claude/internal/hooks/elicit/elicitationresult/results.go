package elicitationresult

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Accept returns an accept action result.
	Accept() Output
	// Decline returns a decline action result that blocks the response.
	Decline() Output
	// Cancel returns a cancel action result.
	Cancel() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Accept returns an accept action result.
func (results) Accept() Output {
	return output{action: "accept"}
}

// Decline returns a decline action result that blocks the response.
func (results) Decline() Output {
	return output{action: "decline"}
}

// Cancel returns a cancel action result.
func (results) Cancel() Output {
	return output{action: "cancel"}
}
