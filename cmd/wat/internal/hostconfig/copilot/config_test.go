package copilot_test

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/copilot"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

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
		Type:           "command",
		Bash:           "echo hi",
		Cwd:            "/tmp",
		Env:            map[string]string{"FOO": "bar"},
		Headers:        map[string]string{"X-Source": "wat"},
		AllowedEnvVars: []string{"GITHUB_TOKEN"},
		TimeoutSec:     30,
	})
	if err != nil {
		t.Fatal(err)
	}
	h, err := hookkit.ParseHandler[copilot.Handler](raw)
	if err != nil {
		t.Fatal(err)
	}
	if h.Bash != "echo hi" || h.TimeoutSec != 30 || h.Cwd != "/tmp" {
		t.Fatalf("handler = %+v", h)
	}
	if h.Env["FOO"] != "bar" || h.Headers["X-Source"] != "wat" || len(h.AllowedEnvVars) != 1 || h.AllowedEnvVars[0] != "GITHUB_TOKEN" {
		t.Fatalf("optional fields = env=%v headers=%v allowedEnvVars=%v", h.Env, h.Headers, h.AllowedEnvVars)
	}
	if h.TimeoutSeconds() != 30 || h.EffectiveCommand() != "echo hi" {
		t.Fatalf("helpers = %d, %q", h.TimeoutSeconds(), h.EffectiveCommand())
	}
}

func TestFile_DisableAllHooksRoundTrip(t *testing.T) {
	t.Parallel()
	in := copilot.File{
		Version:         1,
		DisableAllHooks: true,
		Hooks:           map[string][]json.RawMessage{},
	}
	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out copilot.File
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	if !out.DisableAllHooks || out.Version != 1 {
		t.Fatalf("file = %+v", out)
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
