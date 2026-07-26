package beforetabfileread

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
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

func TestEncode_TabFileReadAllow_permissionOnly(t *testing.T) {
	out, code, err := results{}.Allow().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	assertPermissionOnlyJSON(t, out, "allow")
}

func TestEncode_TabFileReadDeny_permissionOnlyExitZero(t *testing.T) {
	out, code, err := results{}.Deny().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	assertPermissionOnlyJSON(t, out, "deny")
}

func TestEncode_TabFileReadDeny_ignoresChainedMessages(t *testing.T) {
	out, code, err := results{}.Deny().
		WithUserMessage("user").
		WithAgentMessage("agent").
		WithUpdatedInput(map[string]any{"file_path": "x"}).
		Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	assertPermissionOnlyJSON(t, out, "deny")
}

func TestDecode_BeforeTabFileRead(t *testing.T) {
	ev := mustDecode[Event](t, `{"hook_event_name":"beforeTabFileRead","conversation_id":"c1","file_path":"x.go","content":"x"}`)
	if ev.FilePath != "x.go" {
		t.Errorf("FilePath = %q, want x.go", ev.FilePath)
	}
	if ev.Content != "x" {
		t.Errorf("Content = %q, want x", ev.Content)
	}
}

func assertPermissionOnlyJSON(t *testing.T, raw []byte, wantPermission string) {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("json: %v (raw %s)", err, raw)
	}
	if got["permission"] != wantPermission {
		t.Fatalf("permission = %v, want %q (raw %s)", got["permission"], wantPermission, raw)
	}
	if len(got) != 1 {
		t.Fatalf("want permission-only JSON, got %s", raw)
	}
	for _, banned := range []string{"user_message", "agent_message", "updated_input", `"ask"`} {
		if strings.Contains(string(raw), banned) {
			t.Fatalf("beforeTabFileRead must not emit %s: %s", banned, raw)
		}
	}
}

func init() {
	register(testCodec)
}
