package cursor

import (
	"testing"
)

func TestDecode_PreCompact(t *testing.T) {
	mustDecode[PreCompact](t, `{"hook_event_name":"preCompact","conversation_id":"c1","trigger":"auto"}`)
}
