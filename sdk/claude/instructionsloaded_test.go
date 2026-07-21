package claude

import (
	"testing"
)

func TestDecode_InstructionsLoaded(t *testing.T) {
	mustDecode[InstructionsLoaded](t, `{"session_id":"s","hook_event_name":"InstructionsLoaded","file_path":"/f","memory_type":"Project","load_reason":"session_start"}`, EventInstructionsLoaded)
}
