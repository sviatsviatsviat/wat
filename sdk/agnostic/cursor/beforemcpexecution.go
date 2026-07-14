package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapBeforeMCPExecution(e sdkcursor.BeforeMCPExecution, ev *model.Event, name string) {
	nameNorm, _ := model.NormalizeToolName(e.ToolName)
	ev.Tool = &model.ToolCall{
		Name:   nameNorm,
		Native: name,
		Input:  adapter.CloneRaw(e.ToolInput),
		MCP:    true,
	}
}
