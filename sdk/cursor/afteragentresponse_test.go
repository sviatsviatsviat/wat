package cursor

import (
	"testing"
)

func TestDecode_AfterAgentResponse(t *testing.T) {
	mustDecode[AfterAgentResponse](t, `{"hook_event_name":"afterAgentResponse","conversation_id":"c1","text":"done"}`)
}
