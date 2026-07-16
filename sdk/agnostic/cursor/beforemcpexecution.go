package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapBeforeMCPExecution(e sdkcursor.BeforeMCPExecution, ev *model.Event, name string) {
	nameNorm, _ := model.NormalizeToolName(e.ToolName)
	toolInput := model.NewToolInput(nameNorm, e.ToolName, e.ToolInput.Raw())
	ev.Tool = &model.ToolCall{
		Name:   nameNorm,
		Native: name,
		Input:  toolInput,
		MCP:    true,
	}
}
