package stopevent

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

const cursorStop = `{
  "conversation_id": "c1",
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

	out, code, err := results{}.FollowUp("retry with fixed creds").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"followup_message":"retry with fixed creds"`) {
		t.Fatalf("bad stop output: %s", out)
	}
}

func init() {
	register(testCodec)
}
