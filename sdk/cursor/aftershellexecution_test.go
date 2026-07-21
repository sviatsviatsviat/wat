package cursor

import (
	"testing"
)

func TestDecode_AfterShellExecution(t *testing.T) {
	e := mustDecode[AfterShellExecution](t, `{"hook_event_name":"afterShellExecution","conversation_id":"c1","command":"ls","output":"a\nb","duration":10}`)
	if e.Command != "ls" || e.Output != "a\nb" {
		t.Fatalf("event=%+v", e)
	}
}
