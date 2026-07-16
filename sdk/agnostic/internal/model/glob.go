package model

// GlobInput is the canonical glob tool input.
type GlobInput struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path,omitempty"`
}

// AsGlob returns the glob tool input when this payload is for a glob tool.
func (in ToolInput) AsGlob() (GlobInput, bool) {
	return as[GlobInput](in, ToolGlob)
}
