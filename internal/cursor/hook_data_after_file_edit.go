package cursor

import "github.com/sviatsviatsviat/wat/internal/cursor/core"

// hookDataAfterFileEdit is the Cursor payload shape for afterFileEdit.
type hookDataAfterFileEdit = cursorcore.HookDataWithCommon[hookDataAfterFileEditFields]

type hookDataAfterFileEditFields struct {
	FilePath string                          `json:"file_path"`
	Edits    []hookDataAfterFileEditEditPair `json:"edits"`
}

type hookDataAfterFileEditEditPair struct {
	OldString string `json:"old_string"`
	NewString string `json:"new_string"`
}
