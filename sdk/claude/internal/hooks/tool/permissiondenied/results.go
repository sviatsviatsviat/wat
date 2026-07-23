package permissiondenied

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Retry returns a retry-requested PermissionDenied result.
	Retry() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Retry returns a retry-requested PermissionDenied result.
func (results) Retry() Output {
	return output{retry: true}
}
