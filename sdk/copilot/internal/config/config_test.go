package config_test

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/config"
)

func TestHandler_EffectiveCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		h    config.Handler
		want string
	}{
		{name: "command", h: config.Handler{Command: "wat run"}, want: "wat run"},
		{name: "bash", h: config.Handler{Bash: "echo hi"}, want: "echo hi"},
		{name: "powershell", h: config.Handler{PowerShell: "Write-Host hi"}, want: "Write-Host hi"},
		{name: "command precedence", h: config.Handler{Command: "a", Bash: "b", PowerShell: "c"}, want: "a"},
		{name: "bash over powershell", h: config.Handler{Bash: "b", PowerShell: "c"}, want: "b"},
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
	raw, err := hookkit.MarshalHandler(config.Handler{
		Type:       "command",
		Bash:       "echo hi",
		TimeoutSec: 30,
	})
	if err != nil {
		t.Fatal(err)
	}
	h, err := hookkit.ParseHandler[config.Handler](raw)
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
	handlers := []config.Handler{
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
		got, err := hookkit.ParseHandler[config.Handler](blobs[i])
		if err != nil {
			t.Fatalf("blobs[%d]: parse: %v", i, err)
		}
		if got.Command != wantCommand {
			t.Fatalf("blobs[%d].Command = %q, want %q", i, got.Command, wantCommand)
		}
	}
}
