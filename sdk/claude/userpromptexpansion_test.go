package claude

import (
	"testing"
)

func TestDecode_UserPromptExpansion(t *testing.T) {
	mustDecode[UserPromptExpansion](t, `{"session_id":"s","hook_event_name":"UserPromptExpansion","expansion_type":"slash_command","command_name":"foo"}`, EventUserPromptExpansion)
}
