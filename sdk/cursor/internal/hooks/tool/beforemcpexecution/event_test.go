package beforemcpexecution

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() == "" {
		t.Fatal("EventName empty")
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

func TestDecode_BeforeMCPExecution_URLAndCommand(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		wantURL     string
		wantCommand string
		wantTool    string
		wantModelID string
	}{
		{
			name: "url_remote_server",
			raw: `{
				"hook_event_name":"beforeMCPExecution",
				"conversation_id":"c1",
				"generation_id":"g1",
				"model":"gpt-5",
				"model_id":"gpt-5-high",
				"model_params":[{"id":"effort","value":"high"}],
				"cursor_version":"1.7.2",
				"workspace_roots":["/w"],
				"cwd":"/w",
				"tool_name":"MCP:browser_navigate",
				"tool_input":"{\"url\":\"https://example.com\"}",
				"url":"https://mcp.example/mcp"
			}`,
			wantURL:     "https://mcp.example/mcp",
			wantCommand: "",
			wantTool:    "MCP:browser_navigate",
			wantModelID: "gpt-5-high",
		},
		{
			name: "command_stdio_server",
			raw: `{
				"hook_event_name":"beforeMCPExecution",
				"conversation_id":"c1",
				"generation_id":"g1",
				"model":"gpt-5",
				"model_id":"gpt-5-high",
				"cursor_version":"1.7.2",
				"workspace_roots":["/w"],
				"cwd":"/w",
				"tool_name":"MCP:list_resources",
				"tool_input":"{}",
				"command":"npx @example/mcp-server"
			}`,
			wantURL:     "",
			wantCommand: "npx @example/mcp-server",
			wantTool:    "MCP:list_resources",
			wantModelID: "gpt-5-high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := mustDecode[Event](t, tt.raw)
			if e.URL != tt.wantURL {
				t.Fatalf("URL = %q, want %q", e.URL, tt.wantURL)
			}
			if e.Command != tt.wantCommand {
				t.Fatalf("Command = %q, want %q", e.Command, tt.wantCommand)
			}
			if e.ToolName != tt.wantTool {
				t.Fatalf("ToolName = %q, want %q", e.ToolName, tt.wantTool)
			}
			if e.ModelID != tt.wantModelID {
				t.Fatalf("ModelID = %q, want %q", e.ModelID, tt.wantModelID)
			}
			if !e.ToolInput.HasRaw() {
				t.Fatal("ToolInput missing raw payload")
			}
		})
	}
}

func TestEncode_BeforeMCPExecutionDeny(t *testing.T) {
	out, code, err := event.NewPermissionResults().Deny("mcp blocked").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, event.PermissionDenyExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "deny" || got["agent_message"] != "mcp blocked" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_BeforeMCPAsk_enforcedPermission(t *testing.T) {
	out, code, err := event.NewPermissionResults().Ask("confirm MCP tool").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0 (ask is not deny/exit 2)", code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "ask" {
		t.Fatalf("permission = %v, want ask (enforced on beforeMCPExecution)", got["permission"])
	}
	if got["agent_message"] != "confirm MCP tool" {
		t.Fatalf("agent_message = %v, want confirm MCP tool", got["agent_message"])
	}
	if _, ok := got["user_message"]; ok {
		t.Fatalf("Ask default must not set user_message: %s", out)
	}
}

func init() {
	register(testCodec)
}
