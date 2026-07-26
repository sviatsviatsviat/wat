package workspaceopen

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// PluginPaths returns a workspaceOpen result that loads the given absolute
	// plugin directories for the current workspace.
	PluginPaths(paths []string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// PluginPaths returns a workspaceOpen result that loads the given absolute
// plugin directories for the current workspace.
func (results) PluginPaths(paths []string) Output {
	return output{pluginPaths: append([]string(nil), paths...)}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
