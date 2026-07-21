package cursor

import (
	"encoding/json"
	"testing"
)

func TestEncode_SessionStartEnv(t *testing.T) {
	out, code, err := sessionStartResults{}.Noop().
		WithEnv(map[string]string{"K": "V"}).
		WithAdditionalContext("ctx").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		Env map[string]string `json:"env"`
		Ctx string            `json:"additional_context"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.Env["K"] != "V" || got.Ctx != "ctx" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_SessionStart(t *testing.T) {
	e := mustDecode[SessionStart](t, `{"hook_event_name":"sessionStart","conversation_id":"c1","model":"gpt","is_background_agent":true,"cwd":"/w"}`)
	if e.Model != "gpt" || !e.IsBackgroundAgent {
		t.Fatalf("event=%+v", e)
	}
}
