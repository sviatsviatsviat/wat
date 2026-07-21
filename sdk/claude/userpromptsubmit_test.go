package claude

import (
	"strings"
	"testing"
)

func TestEncode_UserPromptBlock(t *testing.T) {
	out, _, err := codec.Encode(userPromptSubmitResults{}.Block("blocked prompt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "blocked prompt") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_UserPromptSubmit(t *testing.T) {
	mustDecode[UserPromptSubmit](t, `{"session_id":"s","hook_event_name":"UserPromptSubmit","prompt":"hello"}`, EventUserPromptSubmit)
}
