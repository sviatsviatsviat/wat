package cursor

import (
	"testing"
)

func TestDecode_AfterAgentThought(t *testing.T) {
	mustDecode[AfterAgentThought](t, `{"hook_event_name":"afterAgentThought","conversation_id":"c1","text":"thinking"}`)
}
