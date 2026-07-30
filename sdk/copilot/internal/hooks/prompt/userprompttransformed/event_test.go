package userprompttransformed

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func init() {
	register(testCodec)
}

func TestDecode_UserPromptTransformed(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"UserPromptTransformed","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","prompt":"hello","transformed_prompt":"HELLO"}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.Prompt != "hello" || e.TransformedPrompt != "HELLO" || e.EventName() != event.UserPromptTransformed {
		t.Fatalf("UserPromptTransformed=%+v", ev)
	}
}

func TestEncode_ModifiedTransformedPrompt(t *testing.T) {
	out, code, err := results{}.Modified("rewritten").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"modified_transformed_prompt":"rewritten"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_Noop(t *testing.T) {
	out, code, err := results{}.Noop().Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if len(out) != 0 {
		t.Fatalf("want empty stdout, got %s", out)
	}
}

func TestMerge_ModifiedTransformedPromptLastWins(t *testing.T) {
	merged, warnings, err := results{}.Modified("first").Merge(results{}.Modified("second"))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) == 0 {
		t.Fatal("want overwrite warning")
	}
	out, code, err := merged.Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"modified_transformed_prompt":"second"`) {
		t.Fatalf("bad merged output: %s", out)
	}
}

func TestRegisterHandler_Nil(t *testing.T) {
	RegisterHandler(nil, nil)
}
