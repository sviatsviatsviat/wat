package checks

// Context carries resolved project state shared across doctor checks.
type Context struct {
	WatDir string
	WatErr error
}
