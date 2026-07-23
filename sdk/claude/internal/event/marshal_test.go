package event

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarshalHookOutput(t *testing.T) {
	t.Run("empty event name", func(t *testing.T) {
		out, code, err := MarshalHookOutput("", func(top, hso map[string]any) {})
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "empty event name") {
			t.Fatalf("error = %v", err)
		}
		if code != SuccessExit {
			t.Fatalf("exit = %d, want %d", code, SuccessExit)
		}
		if out != nil {
			t.Fatalf("out = %q, want nil", out)
		}
	})

	t.Run("empty outputs", func(t *testing.T) {
		out, code, err := MarshalHookOutput(Notification, func(top, hso map[string]any) {})
		if err != nil {
			t.Fatal(err)
		}
		if code != SuccessExit {
			t.Fatalf("exit = %d, want %d", code, SuccessExit)
		}
		if out != nil {
			t.Fatalf("out = %q, want nil", out)
		}
	})

	t.Run("top-level only", func(t *testing.T) {
		out, code, err := MarshalHookOutput(Notification, func(top, hso map[string]any) {
			top["continue"] = false
			top["systemMessage"] = "stop"
		})
		if err != nil {
			t.Fatal(err)
		}
		if code != SuccessExit {
			t.Fatalf("exit = %d, want %d", code, SuccessExit)
		}
		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatal(err)
		}
		if _, ok := got["hookSpecificOutput"]; ok {
			t.Fatalf("unexpected hookSpecificOutput: %s", out)
		}
		if got["continue"] != false || got["systemMessage"] != "stop" {
			t.Fatalf("top-level fields = %s", out)
		}
	})

	t.Run("hook-specific output", func(t *testing.T) {
		out, code, err := MarshalHookOutput(Notification, func(top, hso map[string]any) {
			hso["additionalContext"] = "note"
		})
		if err != nil {
			t.Fatal(err)
		}
		if code != SuccessExit {
			t.Fatalf("exit = %d, want %d", code, SuccessExit)
		}
		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatal(err)
		}
		hso, ok := got["hookSpecificOutput"].(map[string]any)
		if !ok {
			t.Fatalf("missing hookSpecificOutput: %s", out)
		}
		if hso["hookEventName"] != Notification {
			t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], Notification)
		}
		if hso["additionalContext"] != "note" {
			t.Fatalf("additionalContext = %v", hso["additionalContext"])
		}
	})
}
