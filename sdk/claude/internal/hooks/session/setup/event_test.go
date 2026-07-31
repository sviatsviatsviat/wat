package setup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func init() {
	register(testCodec)
}

func TestDecode(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"session_id":"s","hook_event_name":"Setup","trigger":"init"}`))
	if err != nil {
		t.Fatal(err)
	}
	typed, ok := ev.(Event)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if typed.EventName() != event.Setup || typed.Trigger != "init" {
		t.Fatalf("%+v", typed)
	}
}

func TestEncode_Context(t *testing.T) {
	out, code, err := results{}.Context("deps installed").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d", code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	hso, ok := got["hookSpecificOutput"].(map[string]any)
	if !ok {
		t.Fatalf("missing hookSpecificOutput: %s", out)
	}
	if hso["hookEventName"] != event.Setup || hso["additionalContext"] != "deps installed" {
		t.Fatalf("hso = %v", hso)
	}
}

func TestEncode_Env(t *testing.T) {
	envPath := filepath.Join(t.TempDir(), "env.sh")
	t.Setenv("CLAUDE_ENV_FILE", envPath)
	out, code, err := results{}.Noop().WithEnv(map[string]string{"FOO": "bar"}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d", code)
	}
	if out != nil {
		t.Fatalf("env-only result should produce no stdout, got %q", out)
	}
	written, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(written), `export FOO='bar'`) {
		t.Fatalf("env file = %q", written)
	}
}
