package permissionrequest

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

func TestEncode_PermissionRequestDenyInterrupt(t *testing.T) {
	out, code, err := results{}.Deny("blocked").WithInterrupt(true).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.WarnExit {
		t.Fatalf("code=%d, want %d", code, event.WarnExit)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) || !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestAsk(t *testing.T) {
	out, code, err := results{}.Ask("needs user confirmation").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_PermissionRequest(t *testing.T) {
	ev, err := runtime.Codec.Decode([]byte(`{"hook_event_name":"PermissionRequest","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"create","tool_input":{"path":"a.txt"}}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.NativeToolName() != "create" {
		t.Fatalf("PermissionRequest=%+v", ev)
	}
}

func TestMerge_PermissionRequest(t *testing.T) {
	tests := []struct {
		name          string
		a, b          Output
		wantBehavior  string
		wantMessage   string
		wantSuppress  bool
		wantInterrupt bool
	}{
		{
			name:          "hard_deny_beats_allow",
			a:             results{}.Allow(),
			b:             results{}.Deny("blocked"),
			wantBehavior:  "deny",
			wantMessage:   "blocked",
			wantSuppress:  false,
			wantInterrupt: false,
		},
		{
			name:          "ask_soft_deny",
			a:             results{}.Allow(),
			b:             results{}.Ask("confirm"),
			wantBehavior:  "deny",
			wantMessage:   "confirm",
			wantSuppress:  true,
			wantInterrupt: false,
		},
		{
			name:          "hard_deny_beats_ask",
			a:             results{}.Ask("soft"),
			b:             results{}.Deny("hard"),
			wantBehavior:  "deny",
			wantMessage:   "hard",
			wantSuppress:  false,
			wantInterrupt: false,
		},
		{
			name:          "ask_then_ask_keeps_suppress",
			a:             results{}.Ask("a"),
			b:             results{}.Ask("b"),
			wantBehavior:  "deny",
			wantMessage:   "b",
			wantSuppress:  true,
			wantInterrupt: false,
		},
		{
			name:          "interrupt_sticky",
			a:             results{}.Allow().WithInterrupt(true),
			b:             results{}.Deny("no"),
			wantBehavior:  "deny",
			wantMessage:   "no",
			wantSuppress:  false,
			wantInterrupt: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged, warnings, err := tt.a.Merge(tt.b.(hookkit.Output))
			if err != nil {
				t.Fatal(err)
			}
			if len(warnings) != 0 {
				t.Fatalf("warnings = %v", warnings)
			}
			out := merged.(output)
			if out.behavior != tt.wantBehavior || out.message != tt.wantMessage ||
				out.suppressWarnExit != tt.wantSuppress || out.interrupt != tt.wantInterrupt {
				t.Fatalf("got %#v", out)
			}
		})
	}
}

func TestStop_PermissionRequest(t *testing.T) {
	tests := []struct {
		name string
		out  Output
		want bool
	}{
		{name: "hard_deny", out: results{}.Deny("no"), want: true},
		{name: "ask_soft_deny", out: results{}.Ask("confirm"), want: false},
		{name: "allow", out: results{}.Allow(), want: false},
		{name: "noop", out: results{}.Noop(), want: false},
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
