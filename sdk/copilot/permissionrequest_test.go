package copilot

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestEncode_PermissionRequestDenyInterrupt(t *testing.T) {
	out, code, err := permissionRequestResults{}.Deny("blocked").WithInterrupt(true).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != WarnExit {
		t.Fatalf("code=%d, want %d", code, WarnExit)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) || !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestAsk(t *testing.T) {
	out, code, err := permissionRequestResults{}.Ask("needs user confirmation").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_PermissionRequest(t *testing.T) {
	e := mustDecode[PermissionRequest](t, `{"hook_event_name":"PermissionRequest","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"create","tool_input":{"path":"a.txt"}}`, EventPermissionRequest)
	if e.NativeToolName() != "create" {
		t.Fatalf("PermissionRequest=%+v", e)
	}
}

func TestMerge_PermissionRequest(t *testing.T) {
	tests := []struct {
		name          string
		a, b          PermissionRequestOutput
		wantBehavior  string
		wantMessage   string
		wantSuppress  bool
		wantInterrupt bool
	}{
		{
			name:          "hard_deny_beats_allow",
			a:             permissionRequestResults{}.Allow(),
			b:             permissionRequestResults{}.Deny("blocked"),
			wantBehavior:  "deny",
			wantMessage:   "blocked",
			wantSuppress:  false,
			wantInterrupt: false,
		},
		{
			name:          "ask_soft_deny",
			a:             permissionRequestResults{}.Allow(),
			b:             permissionRequestResults{}.Ask("confirm"),
			wantBehavior:  "deny",
			wantMessage:   "confirm",
			wantSuppress:  true,
			wantInterrupt: false,
		},
		{
			name:          "hard_deny_beats_ask",
			a:             permissionRequestResults{}.Ask("soft"),
			b:             permissionRequestResults{}.Deny("hard"),
			wantBehavior:  "deny",
			wantMessage:   "hard",
			wantSuppress:  false,
			wantInterrupt: false,
		},
		{
			name:          "ask_then_ask_keeps_suppress",
			a:             permissionRequestResults{}.Ask("a"),
			b:             permissionRequestResults{}.Ask("b"),
			wantBehavior:  "deny",
			wantMessage:   "b",
			wantSuppress:  true,
			wantInterrupt: false,
		},
		{
			name:          "interrupt_sticky",
			a:             permissionRequestResults{}.Allow().WithInterrupt(true),
			b:             permissionRequestResults{}.Deny("no"),
			wantBehavior:  "deny",
			wantMessage:   "no",
			wantSuppress:  false,
			wantInterrupt: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged, warnings, err := tt.a.Merge(tt.b.(run.Output))
			if err != nil {
				t.Fatal(err)
			}
			if len(warnings) != 0 {
				t.Fatalf("warnings = %v", warnings)
			}
			out := merged.(permissionRequestOutput)
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
		out  PermissionRequestOutput
		want bool
	}{
		{name: "hard_deny", out: permissionRequestResults{}.Deny("no"), want: true},
		{name: "ask_soft_deny", out: permissionRequestResults{}.Ask("confirm"), want: false},
		{name: "allow", out: permissionRequestResults{}.Allow(), want: false},
		{name: "noop", out: permissionRequestResults{}.Noop(), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.out.Stop(); got != tt.want {
				t.Fatalf("Stop() = %v, want %v", got, tt.want)
			}
		})
	}
}
