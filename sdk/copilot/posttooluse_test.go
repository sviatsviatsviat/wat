package copilot

import (
	"strings"
	"testing"
)

func TestEncode_PostToolUpdatedOutput(t *testing.T) {
	out, code, err := postToolResults{}.Context("extra guidance").WithModifiedResult("rewritten").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	s := string(out)
	if !strings.Contains(s, `"text_result_for_llm":"rewritten"`) || !strings.Contains(s, `"additional_context":"extra guidance"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_PostToolUseResultRaw(t *testing.T) {
	raw := `{"hook_event_name":"PostToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"view","tool_input":{},"tool_result":{"result_type":"success","text_result_for_llm":"contents"}}`
	ev, err := codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	post := ev.(PostToolUse)
	got := string(post.ResultRaw())
	if !strings.Contains(got, "text_result_for_llm") || !strings.Contains(got, "contents") {
		t.Fatalf("ResultRaw=%s", got)
	}
}

func TestDecode_PostToolUse(t *testing.T) {
	e := mustDecode[PostToolUse](t, `{"hook_event_name":"PostToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"view","tool_input":{},"tool_result":{"result_type":"success","text_result_for_llm":"contents"}}`, EventPostToolUse)
	if e.ResultText() != "contents" {
		t.Fatalf("PostToolUse=%+v", e)
	}
}
