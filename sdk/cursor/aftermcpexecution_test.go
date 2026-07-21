package cursor

import (
	"testing"
)

func TestDecode_AfterMCPExecution(t *testing.T) {
	mustDecode[AfterMCPExecution](t, `{"hook_event_name":"afterMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","result_json":"{}","duration_ms":5}`)
}
