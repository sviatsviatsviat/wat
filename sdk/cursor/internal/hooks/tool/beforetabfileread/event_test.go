package beforetabfileread

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := runtime.Codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() == "" {
		t.Fatal("EventName empty")
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

func TestEncode_TabFileReadDeny(t *testing.T) {
	out, code, err := event.NewPermissionResults().Deny("no tab reads").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, event.PermissionDenyExit)
	}
	if !strings.Contains(string(out), `"permission":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_BeforeTabFileRead(t *testing.T) {
	mustDecode[Event](t, `{"hook_event_name":"beforeTabFileRead","conversation_id":"c1","file_path":"x.go","content":"x"}`)
}

func init() {
	Register(runtime.Codec)
}
