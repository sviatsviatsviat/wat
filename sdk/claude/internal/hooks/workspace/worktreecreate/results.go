package worktreecreate

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Path returns a worktree-path result.
	Path(path string) Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Path returns a worktree-path result.
func (results) Path(path string) Output {
	return output{worktreePath: path}
}
