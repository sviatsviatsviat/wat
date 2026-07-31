package filechanged

// Results is the hook-scoped response builder for FileChanged.
type Results interface {
	// WatchPaths returns a result that replaces the dynamic FileChanged watch list.
	WatchPaths(paths []string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// WatchPaths returns a result that replaces the dynamic FileChanged watch list.
// A nil or empty slice clears the dynamic watch list (encoded as []).
func (results) WatchPaths(paths []string) Output {
	return output{watchPaths: append([]string{}, paths...)}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
