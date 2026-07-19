package copilot_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

const copilotPreToolUse = `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "bash",
  "tool_input": {"command": "rm -rf /"}
}`

func TestMux_Serve_PreToolHandlerError(t *testing.T) {
	run.Reset()
	copilot.OnPreToolUse(func(ctx context.Context, hook run.Hook[copilot.PreToolUse], _ copilot.PreToolResults) (copilot.PreToolOutput, error) {
		return nil, errors.New("boom")
	})
	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(copilotPreToolUse), &stdout, &bytes.Buffer{})
	if code != copilot.HandlerErrorExit {
		t.Fatalf("exit = %d, want %d", code, copilot.HandlerErrorExit)
	}
}

func TestMux_Serve_PreToolDeny(t *testing.T) {
	run.Reset()
	copilot.OnPreToolUse(func(ctx context.Context, hook run.Hook[copilot.PreToolUse], r copilot.PreToolResults) (copilot.PreToolOutput, error) {
		return r.Deny("nope"), nil
	})
	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(copilotPreToolUse), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stdout.String(), `"permission_decision":"deny"`) {
		t.Fatalf("stdout=%s", stdout.String())
	}
}

func TestHandler_EffectiveCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		h    copilot.Handler
		want string
	}{
		{name: "command", h: copilot.Handler{Command: "wat run"}, want: "wat run"},
		{name: "bash", h: copilot.Handler{Bash: "echo hi"}, want: "echo hi"},
		{name: "powershell", h: copilot.Handler{PowerShell: "Write-Host hi"}, want: "Write-Host hi"},
		{name: "command precedence", h: copilot.Handler{Command: "a", Bash: "b", PowerShell: "c"}, want: "a"},
		{name: "bash over powershell", h: copilot.Handler{Bash: "b", PowerShell: "c"}, want: "b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.h.EffectiveCommand(); got != tt.want {
				t.Fatalf("EffectiveCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseHandler_RoundTrip(t *testing.T) {
	raw, err := hookkit.MarshalHandler(copilot.Handler{
		Type:       "command",
		Bash:       "echo hi",
		TimeoutSec: 30,
	})
	if err != nil {
		t.Fatal(err)
	}
	h, err := hookkit.ParseHandler[copilot.Handler](raw)
	if err != nil {
		t.Fatal(err)
	}
	if h.Bash != "echo hi" || h.TimeoutSec != 30 {
		t.Fatalf("handler = %+v", h)
	}
	if h.TimeoutSeconds() != 30 || h.EffectiveCommand() != "echo hi" {
		t.Fatalf("helpers = %d, %q", h.TimeoutSeconds(), h.EffectiveCommand())
	}
}

func TestHandlers_EncodesMultiple(t *testing.T) {
	handlers := []copilot.Handler{
		{Type: "command", Command: "a"},
		{Type: "command", Command: "b"},
	}
	blobs := make([]json.RawMessage, 0, len(handlers))
	for _, h := range handlers {
		raw, err := hookkit.MarshalHandler(h)
		if err != nil {
			t.Fatal(err)
		}
		blobs = append(blobs, raw)
	}
	if len(blobs) != 2 {
		t.Fatalf("len = %d", len(blobs))
	}
	for i, wantCommand := range []string{"a", "b"} {
		got, err := hookkit.ParseHandler[copilot.Handler](blobs[i])
		if err != nil {
			t.Fatalf("blobs[%d]: parse: %v", i, err)
		}
		if got.Command != wantCommand {
			t.Fatalf("blobs[%d].Command = %q, want %q", i, got.Command, wantCommand)
		}
	}
}
