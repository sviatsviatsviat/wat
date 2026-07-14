package claude

// FileChanged is the FileChanged hook event.
type FileChanged struct {
	Envelope
	// FilePath is the changed file path.
	FilePath string `json:"file_path"`
}

// EventName returns the hook event name.
func (FileChanged) EventName() string { return EventFileChanged }

func init() {
	registerDecoder(EventFileChanged, decodeAs[FileChanged])
}
