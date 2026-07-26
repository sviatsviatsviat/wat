package stopevent

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

const cursorStop = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "claude-opus-4-7-thinking-max",
  "hook_event_name": "stop",
  "cursor_version": "1.7.2",
  "status": "error",
  "loop_count": 1
}`

func TestDecodeEncode_StopFollowUp(t *testing.T) {
	ev, err := testCodec.Decode([]byte(cursorStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(Event)
	if !ok {
		t.Fatalf("want Event, got %T", ev)
	}
	if stop.Status != "error" || stop.LoopCount != 1 {
		t.Fatalf("bad stop: %+v", stop)
	}
	if stop.Model != "claude-opus-4-7-thinking-max" {
		t.Fatalf("Model = %q", stop.Model)
	}

	out, code, err := results{}.FollowUp("retry with fixed creds").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var payload map[string]any
	if err := json.Unmarshal(out, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["followup_message"] != "retry with fixed creds" {
		t.Fatalf("bad stop output: %s", out)
	}
	if len(payload) != 1 {
		t.Fatalf("want only followup_message, got %s", out)
	}
}

func TestEncode_StopNoopSilent(t *testing.T) {
	out, code, err := results{}.Noop().Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if len(out) != 0 {
		t.Fatalf("noop should be silent, got %q", out)
	}
}

func TestEncode_StopEmptyFollowUpSilent(t *testing.T) {
	out, code, err := results{}.FollowUp("").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if len(out) != 0 {
		t.Fatalf("empty follow-up should be silent, got %q", out)
	}
}

func TestMerge_StopFollowUpTakeLast(t *testing.T) {
	merged, warnings, err := results{}.FollowUp("first").Merge(results{}.FollowUp("second"))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) == 0 || !strings.Contains(warnings[0], "followUpMessage") {
		t.Fatalf("want take-last warning, got %v", warnings)
	}
	out, code, err := merged.Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"followup_message":"second"`) {
		t.Fatalf("merged output = %s", out)
	}
}

func init() {
	register(testCodec)
}
