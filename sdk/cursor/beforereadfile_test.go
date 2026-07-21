package cursor

import (
	"testing"
)

func TestDecode_BeforeReadFile(t *testing.T) {
	mustDecode[BeforeReadFile](t, `{"hook_event_name":"beforeReadFile","conversation_id":"c1","file_path":"a.go","content":"package main"}`)
}
