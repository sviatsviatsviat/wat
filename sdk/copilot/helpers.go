package copilot

import (
	"encoding/json"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
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

func extractRawObjectField(raw json.RawMessage, camelKey, snakeKey string) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		return nil
	}
	if b, ok := fields[camelKey]; ok && len(b) > 0 && string(b) != "null" {
		return hookkit.CloneRaw(b)
	}
	if b, ok := fields[snakeKey]; ok && len(b) > 0 && string(b) != "null" {
		return hookkit.CloneRaw(b)
	}
	return nil
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
