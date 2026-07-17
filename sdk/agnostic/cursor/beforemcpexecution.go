package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapBeforeMCPExecution maps a Cursor BeforeMCPExecution hook into a unified Event.
func MapBeforeMCPExecution(e sdkcursor.BeforeMCPExecution, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreTool)
	name := receivedName(e)
	nameNorm, _ := model.NormalizeToolName(e.ToolName)
	toolInput := model.NewToolInput(nameNorm, e.ToolName, e.ToolInput.Raw())
	ev.Tool = &model.ToolCall{
		Name:   nameNorm,
		Native: name,
		Input:  toolInput,
		MCP:    true,
	}
	return ev
}
