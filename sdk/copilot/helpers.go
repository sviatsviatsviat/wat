package copilot

import (
	"encoding/json"
	"strings"

	"github.com/sviatsviatsviat/wat/sdk/copilot/tools"
)

func isShellToolName(name string) bool {
	switch strings.ToLower(name) {
	case tools.ToolBash, tools.ToolPowerShell, tools.ToolShell:
		return true
	default:
		return false
	}
}

func marshalToolResultCamel(r ToolResult) json.RawMessage {
	out := map[string]string{}
	if r.ResultType != "" {
		out["resultType"] = r.ResultType
	}
	if r.TextResultForLLM != "" {
		out["textResultForLlm"] = r.TextResultForLLM
	}
	if len(out) == 0 {
		return nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil
	}
	return b
}

func marshalToolResultSnake(r ToolResult) json.RawMessage {
	out := map[string]string{}
	if r.ResultTypeSnake != "" {
		out["result_type"] = r.ResultTypeSnake
	}
	if r.TextResultSnake != "" {
		out["text_result_for_llm"] = r.TextResultSnake
	}
	if len(out) == 0 {
		return nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil
	}
	return b
}
