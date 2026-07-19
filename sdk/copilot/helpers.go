package copilot

import (
	"encoding/json"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func isShellToolName(name string) bool {
	switch strings.ToLower(name) {
	case tools.ToolBash, tools.ToolPowerShell, tools.ToolShell:
		return true
	default:
		return false
	}
}

func marshalToolResult(r ToolResult) json.RawMessage {
	out := map[string]string{}
	if r.ResultType != "" {
		out["result_type"] = r.ResultType
	}
	if r.TextResultForLLM != "" {
		out["text_result_for_llm"] = r.TextResultForLLM
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

type nativeToolNamer interface {
	NativeToolName() string
}

// registerToolInputEvent registers a decoder that populates ToolInput from "tool_input".
func registerToolInputEvent[T hookkit.Event, PT interface {
	*T
	nativeToolNamer
}](eventName string, toolInput func(*T) *tools.Input) {
	codec.Register(eventName, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *T, payload []byte) {
			*toolInput(e) = tools.NewInputFromPayload(PT(e).NativeToolName(), payload, "tool_input")
		})
	})
}
