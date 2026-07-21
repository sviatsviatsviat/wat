package cursor

import (
	"testing"
)

func TestDecode_BeforeMCPExecution(t *testing.T) {
	e := mustDecode[BeforeMCPExecution](t, `{"hook_event_name":"beforeMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":"{}","url":"https://mcp.example/mcp"}`)
	if e.URL != "https://mcp.example/mcp" {
		t.Fatalf("URL=%q", e.URL)
	}
}
