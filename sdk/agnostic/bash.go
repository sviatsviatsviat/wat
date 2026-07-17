package agnostic

// BashInput is the canonical shell tool input.
type BashInput struct {
	Command string `json:"command"`
}

// AsBash returns the bash tool input when this payload is for a shell tool.
func (in ToolInput) AsBash() (BashInput, bool) {
	return as[BashInput](in, ToolBash)
}
