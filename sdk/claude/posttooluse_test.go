package claude

import (
	"encoding/json"
	"testing"
)

func TestDecode_PostToolUse(t *testing.T) {
	mustDecode[PostToolUse](t, `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`, EventPostToolUse)

	t.Run("context", func(t *testing.T) {
		out, code, err := postToolUseResults{}.Context("tool finished").Encode()
		if err != nil {
			t.Fatal(err)
		}
		if code != SuccessExit {
			t.Fatalf("exit = %d, want %d", code, SuccessExit)
		}
		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatal(err)
		}
		hso, ok := got["hookSpecificOutput"].(map[string]any)
		if !ok {
			t.Fatalf("missing hookSpecificOutput: %s", out)
		}
		if hso["hookEventName"] != EventPostToolUse {
			t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], EventPostToolUse)
		}
		if hso["additionalContext"] != "tool finished" {
			t.Fatalf("additionalContext = %v", hso["additionalContext"])
		}
	})

	t.Run("block", func(t *testing.T) {
		out, code, err := postToolUseResults{}.Block("unsafe output").Encode()
		if err != nil {
			t.Fatal(err)
		}
		if code != SuccessExit {
			t.Fatalf("exit = %d, want %d", code, SuccessExit)
		}
		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatal(err)
		}
		if got["decision"] != "block" || got["reason"] != "unsafe output" {
			t.Fatalf("block fields = %s", out)
		}
		// Block-only responses put decision/reason at top level; no hookSpecificOutput.
		if _, ok := got["hookSpecificOutput"]; ok {
			t.Fatalf("unexpected hookSpecificOutput: %s", out)
		}
	})
}
