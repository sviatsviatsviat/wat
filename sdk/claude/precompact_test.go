package claude

import (
	"testing"
)

func TestDecode_PreCompact(t *testing.T) {
	mustDecode[PreCompact](t, `{"session_id":"s","hook_event_name":"PreCompact","trigger":"manual","custom_instructions":"keep tests"}`, EventPreCompact)
}
