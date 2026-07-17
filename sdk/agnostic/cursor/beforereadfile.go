package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapBeforeReadFile maps a Cursor BeforeReadFile hook into a unified Event.
func MapBeforeReadFile(e sdkcursor.BeforeReadFile, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreTool)
	name := receivedName(e)
	ev.Tool = &model.ToolCall{Name: model.ToolRead, Native: name}
	input, err := json.Marshal(map[string]any{
		"file_path":   e.FilePath,
		"content":     e.Content,
		"attachments": e.Attachments,
	})
	if err != nil {
		return ev
	}
	ev.Tool.Input = model.NewToolInput(model.ToolRead, name, input)
	return ev
}
