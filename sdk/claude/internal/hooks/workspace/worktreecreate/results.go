package worktreecreate

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Path returns a worktree-path result for command-hook stdout.
	Path(path string) Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Path returns a worktree-path result.
// Encode prints path as plain stdout text for Claude Code command hooks; do
// not wrap it in JSON. HTTP hooks use hookSpecificOutput.worktreePath instead.
func (results) Path(path string) Output {
	return output{worktreePath: path}
}
