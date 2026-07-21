package claude

import (
	"testing"
)

func TestDecode_Elicitation(t *testing.T) {
	mustDecode[Elicitation](t, `{"session_id":"s","hook_event_name":"Elicitation","server_name":"srv","message":"confirm?"}`, EventElicitation)
}
