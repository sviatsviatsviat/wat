package beforesubmitprompt

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/stopevent"
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

func TestEncode_BeforeSubmitPromptBlock(t *testing.T) {
	out, code, err := results{}.Block("blocked").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["continue"] != false || got["user_message"] != "blocked" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_BeforeSubmitPrompt(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"beforeSubmitPrompt","conversation_id":"c1","prompt":"hello"}`)
	if e.Prompt != "hello" {
		t.Fatal("bad prompt")
	}
}

func TestMerge_BeforeSubmitPrompt(t *testing.T) {
	t.Run("continue_false_sticky", func(t *testing.T) {
		a := results{}.Block("first")
		b := output{}.WithContinue(true).WithUserMessage("second")
		merged, warnings, err := a.Merge(b.(hookkit.Output))
		if err != nil {
			t.Fatal(err)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], "userMessage") {
			t.Fatalf("warnings = %v", warnings)
		}
		out := merged.(output)
		if out.cont == nil || *out.cont || out.userMessage != "second" {
			t.Fatalf("got %#v", out)
		}
		if !merged.Stop() {
			t.Fatal("continue false should stop")
		}
	})

	t.Run("later_continue_true", func(t *testing.T) {
		a := results{}.Noop()
		b := output{}.WithContinue(true)
		merged, warnings, err := a.Merge(b.(hookkit.Output))
		if err != nil {
			t.Fatal(err)
		}
		if len(warnings) != 0 {
			t.Fatalf("warnings = %v", warnings)
		}
		out := merged.(output)
		if out.cont == nil || !*out.cont {
			t.Fatalf("got %#v", out)
		}
		if merged.Stop() {
			t.Fatal("continue true should not stop")
		}
	})

	t.Run("type_mismatch", func(t *testing.T) {
		_, _, err := results{}.Noop().Merge(stopevent.NewResults().Noop().(hookkit.Output))
		if err == nil || !strings.Contains(err.Error(), "merge type mismatch") {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestStop_BeforeSubmitPrompt(t *testing.T) {
	tests := []struct {
		name string
		out  output
		want bool
	}{
		{name: "nil", out: output{}, want: false},
		{name: "true", out: output{}.WithContinue(true).(output), want: false},
		{name: "false", out: results{}.Block("x").(output), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.out.Stop(); got != tt.want {
				t.Fatalf("Stop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func init() {
	Register(runtime.Codec)
}
