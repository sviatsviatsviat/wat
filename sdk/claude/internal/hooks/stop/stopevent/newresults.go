package stopevent

// NewResults returns a Results builder stamped with eventName.
func NewResults(eventName string) Results {
	return results{eventName: eventName}
}
