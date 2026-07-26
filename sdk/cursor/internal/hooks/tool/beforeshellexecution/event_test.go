package beforeshellexecution

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

const cursorShell = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "some-model",
  "hook_event_name": "beforeShellExecution",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "user_email": null,
  "transcript_path": null,
  "command": "git push --force",
  "cwd": "/w",
  "sandbox": false
}`

func TestDecodeEncode_BeforeShellDeny(t *testing.T) {
	ev, err := testCodec.Decode([]byte(cursorShell))
	if err != nil {
		t.Fatal(err)
	}
	shell, ok := ev.(Event)
	if !ok {
		t.Fatalf("want BeforeShellExecution, got %T", ev)
	}
	if shell.Command != "git push --force" || shell.ConversationID != "c1" {
		t.Fatalf("bad event: %+v", shell)
	}

	out, code, err := event.NewPermissionResults().Deny("force push blocked").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, event.PermissionDenyExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "deny" || got["agent_message"] != "force push blocked" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_BeforeShellAsk_enforcedPermission(t *testing.T) {
	out, code, err := event.NewPermissionResults().Ask("confirm force push").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0 (ask is not deny/exit 2)", code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "ask" {
		t.Fatalf("permission = %v, want ask (enforced on beforeShellExecution)", got["permission"])
	}
	if got["agent_message"] != "confirm force push" {
		t.Fatalf("agent_message = %v, want confirm force push", got["agent_message"])
	}
	if _, ok := got["user_message"]; ok {
		t.Fatalf("Ask default must not set user_message: %s", out)
	}
}

func init() {
	register(testCodec)
}
