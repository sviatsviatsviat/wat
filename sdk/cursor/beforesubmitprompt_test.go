package cursor

import (
	"encoding/json"
	"testing"
)

func TestEncode_BeforeSubmitPromptBlock(t *testing.T) {
	out, code, err := beforeSubmitPromptResults{}.Block("blocked").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["continue"] != false || got["user_message"] != "blocked" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_BeforeSubmitPrompt(t *testing.T) {
	e := mustDecode[BeforeSubmitPrompt](t, `{"hook_event_name":"beforeSubmitPrompt","conversation_id":"c1","prompt":"hello"}`)
	if e.Prompt != "hello" {
		t.Fatal("bad prompt")
	}
}
