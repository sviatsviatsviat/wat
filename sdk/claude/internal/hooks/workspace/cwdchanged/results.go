package cwdchanged

// Results is the hook-scoped response builder for CwdChanged.
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
// An empty slice clears the dynamic watch list.
func (results) WatchPaths(paths []string) Output {
	if paths == nil {
		return output{watchPaths: []string{}}
	}
	return output{watchPaths: append([]string(nil), paths...)}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
