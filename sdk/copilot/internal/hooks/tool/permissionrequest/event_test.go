package permissionrequest

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

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

func TestEncode_PermissionRequestAsk_softDenyNotUserPrompt(t *testing.T) {
	// Ask must not encode a confirmation / "ask" decision: Copilot's
	// permissionRequest schema is allow|deny only. Soft deny is behavior deny
	// with exit 0; Noop/empty stdout is the fall-through to user prompting.
	out, code, err := results{}.Ask("needs user confirmation").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("Ask soft deny code=%d, want 0 (not WarnExit)", code)
	}
	got := string(out)
	if !strings.Contains(got, `"behavior":"deny"`) {
		t.Fatalf("Ask must encode soft deny: %s", got)
	}
	if strings.Contains(got, `"ask"`) {
		t.Fatalf("Ask must not encode ask on the wire: %s", got)
	}
	if !strings.Contains(got, `"message":"needs user confirmation"`) {
		t.Fatalf("Ask must keep message: %s", got)
	}
	if (results{}).Ask("x").Stop() {
		t.Fatal("Ask soft deny must not Stop remaining handlers")
	}

	denyOut, denyCode, err := results{}.Deny("blocked").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if denyCode != event.WarnExit {
		t.Fatalf("Deny code=%d, want WarnExit %d", denyCode, event.WarnExit)
	}
	if !strings.Contains(string(denyOut), `"behavior":"deny"`) {
		t.Fatalf("Deny must encode deny: %s", denyOut)
	}
	if !(results{}).Deny("blocked").Stop() {
		t.Fatal("hard Deny must Stop remaining handlers")
	}

	noopOut, noopCode, err := results{}.Noop().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if noopCode != 0 || len(noopOut) != 0 {
		t.Fatalf("Noop must be empty stdout exit 0, got code=%d out=%q", noopCode, noopOut)
	}
}

func TestDecode_PermissionRequest(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"PermissionRequest","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"create","tool_input":{"path":"a.txt"}}`))
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
	register(testCodec)
}
