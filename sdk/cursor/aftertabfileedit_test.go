package cursor

import (
	"testing"
)

func TestDecode_AfterTabFileEdit(t *testing.T) {
	mustDecode[AfterTabFileEdit](t, `{"hook_event_name":"afterTabFileEdit","conversation_id":"c1","file_path":"x.go","edits":[]}`)
}
