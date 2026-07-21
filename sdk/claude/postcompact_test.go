package claude

import (
	"testing"
)

func TestDecode_PostCompact(t *testing.T) {
	mustDecode[PostCompact](t, `{"session_id":"s","hook_event_name":"PostCompact","trigger":"auto"}`, EventPostCompact)
}
