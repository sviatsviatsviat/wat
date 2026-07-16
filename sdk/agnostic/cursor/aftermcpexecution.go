package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapAfterMCPExecution(e sdkcursor.AfterMCPExecution, ev *model.Event, name string) {
	nameNorm, _ := model.NormalizeToolName(e.ToolName)
	toolInput := model.NewToolInput(nameNorm, e.ToolName, e.ToolInput.Raw())
	ev.Tool = &model.ToolCall{
		Name:   nameNorm,
		Native: name,
		Input:  toolInput,
		MCP:    true,
	}
	ev.Result = &model.ToolResult{Raw: json.RawMessage(e.ResultJSON), DurationMs: e.DurationMillis()}
}
