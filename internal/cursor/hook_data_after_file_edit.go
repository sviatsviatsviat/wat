package cursor

type hookDataAfterFileEditFields struct {
	FilePath string                          `json:"file_path"`
	Edits    []hookDataAfterFileEditEditPair `json:"edits"`
}

type hookDataAfterFileEditEditPair struct {
	OldString string `json:"old_string"`
	NewString string `json:"new_string"`
}
