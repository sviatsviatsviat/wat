package event

// Attachment is a file or rule attachment on prompt and read-file events.
type Attachment struct {
	// Type is the attachment type.
	Type string `json:"type"`
	// FilePath is the attached file path.
	FilePath string `json:"file_path"`
}

// Edit is a single file edit on afterFileEdit events.
type Edit struct {
	// OldString is the text replaced.
	OldString string `json:"old_string"`
	// NewString is the replacement text.
	NewString string `json:"new_string"`
}

// EditRange is the character range of a Tab edit on afterTabFileEdit events.
type EditRange struct {
	// StartLineNumber is the start line of the edit.
	StartLineNumber int `json:"start_line_number"`
	// StartColumn is the start column of the edit.
	StartColumn int `json:"start_column"`
	// EndLineNumber is the end line of the edit.
	EndLineNumber int `json:"end_line_number"`
	// EndColumn is the end column of the edit.
	EndColumn int `json:"end_column"`
}

// TabEdit is a single file edit on afterTabFileEdit events.
//
// Unlike [Edit], Tab edits include precise range and line context from Cursor's
// Tab completion hooks.
type TabEdit struct {
	// OldString is the text replaced.
	OldString string `json:"old_string"`
	// NewString is the replacement text.
	NewString string `json:"new_string"`
	// Range is the character range of the edit.
	Range EditRange `json:"range"`
	// OldLine is the full line content before the edit.
	OldLine string `json:"old_line"`
	// NewLine is the full line content after the edit.
	NewLine string `json:"new_line"`
}
