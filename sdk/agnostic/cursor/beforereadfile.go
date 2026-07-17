package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapBeforeReadFile(e sdkcursor.BeforeReadFile, ev *model.Event, name string) {
	ev.Tool = &model.ToolCall{Name: model.ToolRead, Native: name}
	input, err := json.Marshal(map[string]any{
		"file_path":   e.FilePath,
		"content":     e.Content,
		"attachments": e.Attachments,
	})
	if err != nil {
		return
	}
	ev.Tool.Input = model.NewToolInput(model.ToolRead, name, input)
}
