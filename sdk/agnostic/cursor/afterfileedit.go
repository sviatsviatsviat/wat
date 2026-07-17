package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapAfterFileEdit maps a Cursor AfterFileEdit hook into a unified Event.
func MapAfterFileEdit(e sdkcursor.AfterFileEdit, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPostTool)
	name := receivedName(e)
	ev.Tool = &model.ToolCall{Name: model.ToolEdit, Native: name}
	input, err := json.Marshal(map[string]any{
		"file_path": e.FilePath,
		"edits":     e.Edits,
	})
	if err != nil {
		return ev
	}
	editsRaw, err := json.Marshal(e.Edits)
	if err != nil {
		return ev
	}
	ev.Tool.Input = model.NewToolInput(model.ToolEdit, name, input)
	ev.Result = &model.ToolResult{Raw: editsRaw}
	return ev
}
