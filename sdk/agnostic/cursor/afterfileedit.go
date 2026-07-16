package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapAfterFileEdit(e sdkcursor.AfterFileEdit, ev *model.Event, name string) {
	input, err := json.Marshal(map[string]any{
		"file_path": e.FilePath,
		"edits":     e.Edits,
	})
	if err != nil {
		return
	}
	editsRaw, err := json.Marshal(e.Edits)
	if err != nil {
		return
	}
	toolInput := model.NewToolInput(model.ToolEdit, name, input)
	ev.Tool = &model.ToolCall{Name: model.ToolEdit, Native: name, Input: toolInput}
	ev.Result = &model.ToolResult{Raw: editsRaw}
}
