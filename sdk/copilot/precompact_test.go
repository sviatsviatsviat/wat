package copilot

import (
	"testing"
)

func TestDecode_PreCompact(t *testing.T) {
	e := mustDecode[PreCompact](t, `{"hook_event_name":"PreCompact","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","trigger":"auto","custom_instructions":"keep"}`, EventPreCompact)
	if e.Instructions() != "keep" {
		t.Fatalf("PreCompact=%+v", e)
	}
}
