package cursor

import (
	"strings"
	"testing"
)

const cursorStop = `{
  "conversation_id": "c1",
  "hook_event_name": "stop",
  "cursor_version": "1.7.2",
  "status": "error",
  "loop_count": 1
}`

func TestDecodeEncode_StopFollowUp(t *testing.T) {
	ev, err := codec.Decode([]byte(cursorStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(Stop)
	if !ok {
		t.Fatalf("want Stop, got %T", ev)
	}
	if stop.Status != "error" || stop.LoopCount != 1 {
		t.Fatalf("bad stop: %+v", stop)
	}

	out, code, err := codec.Encode(stopResults{}.FollowUp("retry with fixed creds"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"followup_message":"retry with fixed creds"`) {
		t.Fatalf("bad stop output: %s", out)
	}
}
