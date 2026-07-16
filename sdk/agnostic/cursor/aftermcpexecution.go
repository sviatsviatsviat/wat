package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapAfterMCPExecution(e sdkcursor.AfterMCPExecution, ev *model.Event, name string) {
	nameNorm, _ := model.NormalizeToolName(e.ToolName)
	ev.Tool = &model.ToolCall{
		Name:   nameNorm,
		Native: name,
		Input:  adapter.CloneRaw(e.ToolInput),
		MCP:    true,
	}
	ev.Result = &model.ToolResult{Raw: json.RawMessage(e.ResultJSON), DurationMs: e.DurationMillis()}
}
