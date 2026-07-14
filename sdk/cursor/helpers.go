package cursor

// Attachment is a file or rule attachment on prompt and read-file events.
type Attachment struct {
	// Type is the attachment type.
	Type string `json:"type"`
	// FilePath is the attached file path.
	FilePath string `json:"file_path"`
}

// Edit is a single file edit on afterFileEdit and afterTabFileEdit events.
type Edit struct {
	// OldString is the text replaced.
	OldString string `json:"old_string"`
	// NewString is the replacement text.
	NewString string `json:"new_string"`
}
