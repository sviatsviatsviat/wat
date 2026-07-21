package claude

import (
	"testing"
)

func TestDecode_PostToolBatch(t *testing.T) {
	mustDecode[PostToolBatch](t, `{"session_id":"s","hook_event_name":"PostToolBatch"}`, EventPostToolBatch)
}
