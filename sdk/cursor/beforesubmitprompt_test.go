package cursor

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestEncode_BeforeSubmitPromptBlock(t *testing.T) {
	out, code, err := beforeSubmitPromptResults{}.Block("blocked").Encode()
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
	e := mustDecode[BeforeSubmitPrompt](t, `{"hook_event_name":"beforeSubmitPrompt","conversation_id":"c1","prompt":"hello"}`)
	if e.Prompt != "hello" {
		t.Fatal("bad prompt")
	}
}

func TestMerge_BeforeSubmitPrompt(t *testing.T) {
	t.Run("continue_false_sticky", func(t *testing.T) {
		a := beforeSubmitPromptResults{}.Block("first")
		b := beforeSubmitPromptOutput{}.WithContinue(true).WithUserMessage("second")
		merged, warnings, err := a.Merge(b.(run.Output))
		if err != nil {
			t.Fatal(err)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], "userMessage") {
			t.Fatalf("warnings = %v", warnings)
		}
		out := merged.(beforeSubmitPromptOutput)
		if out.cont == nil || *out.cont || out.userMessage != "second" {
			t.Fatalf("got %#v", out)
		}
		if !merged.Stop() {
			t.Fatal("continue false should stop")
		}
	})

	t.Run("later_continue_true", func(t *testing.T) {
		a := beforeSubmitPromptResults{}.Noop()
		b := beforeSubmitPromptOutput{}.WithContinue(true)
		merged, warnings, err := a.Merge(b.(run.Output))
		if err != nil {
			t.Fatal(err)
		}
		if len(warnings) != 0 {
			t.Fatalf("warnings = %v", warnings)
		}
		out := merged.(beforeSubmitPromptOutput)
		if out.cont == nil || !*out.cont {
			t.Fatalf("got %#v", out)
		}
		if merged.Stop() {
			t.Fatal("continue true should not stop")
		}
	})

	t.Run("type_mismatch", func(t *testing.T) {
		_, _, err := beforeSubmitPromptResults{}.Noop().Merge(stopResults{}.Noop().(run.Output))
		if err == nil || !strings.Contains(err.Error(), "merge type mismatch") {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestStop_BeforeSubmitPrompt(t *testing.T) {
	tests := []struct {
		name string
		out  beforeSubmitPromptOutput
		want bool
	}{
		{name: "nil", out: beforeSubmitPromptOutput{}, want: false},
		{name: "true", out: beforeSubmitPromptOutput{}.WithContinue(true).(beforeSubmitPromptOutput), want: false},
		{name: "false", out: beforeSubmitPromptResults{}.Block("x").(beforeSubmitPromptOutput), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.out.Stop(); got != tt.want {
				t.Fatalf("Stop() = %v, want %v", got, tt.want)
			}
		})
	}
}
