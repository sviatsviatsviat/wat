package cursor

import "testing"

func TestToolInput_AsShell(t *testing.T) {
	ev, err := codec.Decode([]byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(PreToolUse)
	input, ok := pre.ToolInput.AsShell()
	if !ok || input.Command != "ls" {
		t.Fatalf("AsShell = %+v, %v", input, ok)
	}
}

func TestDecode_PreToolUse(t *testing.T) {
	e := mustDecode[PreToolUse](t, `{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`)
	if e.ShellCommand() != "ls" {
		t.Fatalf("ShellCommand=%q", e.ShellCommand())
	}
}
