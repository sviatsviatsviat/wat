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

const copilotCamelPreToolUse = `{
  "sessionId": "s1",
  "timestamp": 1760000000000,
  "cwd": "/w",
  "toolName": "bash",
  "toolArgs": {"command": "rm -rf /"}
}`

const copilotVSCodeStop = `{
  "hook_event_name": "Stop",
  "session_id": "s2",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "transcript_path": "/tmp/t",
  "stop_reason": "end_turn"
}`

func TestDecodeEncode_PreToolDeny(t *testing.T) {
	ev, err := copilot.Decode([]byte(copilotCamelPreToolUse), copilot.WithEvent("preToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(copilot.PreToolUse)
	if !ok || pre.Session() != "s1" || pre.Cwd != "/w" {
		t.Fatalf("bad event: %+v", ev)
	}
	if pre.NativeToolName() != "bash" || pre.ShellCommand() != "rm -rf /" {
		t.Fatalf("bad tool: name=%q shell=%q", pre.NativeToolName(), pre.ShellCommand())
	}
	if !bytes.Equal(copilot.RawBytes(ev), []byte(copilotCamelPreToolUse)) {
		t.Fatal("Raw not preserved")
	}

	out, code, err := copilot.Encode(copilot.EventPreToolUse, copilot.PreToolResultsForTest().Deny("destructive command"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		Decision string `json:"permissionDecision"`
		Reason   string `json:"permissionDecisionReason"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.Decision != "deny" || got.Reason != "destructive command" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_CamelCaseRequiresWithEvent(t *testing.T) {
	_, err := copilot.Decode([]byte(copilotCamelPreToolUse))
	if err == nil {
		t.Fatal("expected error without WithEvent")
	}
	if !strings.Contains(err.Error(), "WithEvent") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_InvalidTypedEvent(t *testing.T) {
	raw := []byte(`{"hook_event_name":"PreToolUse","session_id":"s1","cwd":"/w","tool_name":123}`)
	_, err := copilot.Decode(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, copilot.ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func TestDecode_VSCodeStop(t *testing.T) {
	ev, err := copilot.Decode([]byte(copilotVSCodeStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(copilot.AgentStop)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if stop.EventName() != copilot.EventAgentStop {
		t.Fatalf("EventName=%q", stop.EventName())
	}
	if stop.ReceivedName() != "Stop" {
		t.Fatalf("ReceivedName=%q", stop.ReceivedName())
	}
	if stop.Reason() != "end_turn" {
		t.Fatalf("Reason=%q", stop.Reason())
	}
	if stop.Transcript() != "/tmp/t" {
		t.Fatalf("Transcript=%q", stop.Transcript())
	}
}

func TestDecode_VSCodePreToolBash(t *testing.T) {
	raw := `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "Bash",
  "tool_input": {"command": "ls -la"}
}`
	ev, err := copilot.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(copilot.PreToolUse)
	if !ok || pre.NativeToolName() != "Bash" || pre.ShellCommand() != "ls -la" {
		t.Fatalf("PreToolUse=%+v", ev)
	}
}

func TestDecode_NotificationWithoutWithEvent(t *testing.T) {
	raw := `{
  "sessionId": "s1",
  "timestamp": 1760000000000,
  "cwd": "/w",
  "hook_event_name": "Notification",
  "message": "shell done",
  "title": "Shell completed",
  "notification_type": "shell_completed"
}`
	ev, err := copilot.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	note, ok := ev.(copilot.Notification)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if note.EventName() != copilot.EventNotification {
		t.Fatalf("EventName=%q", note.EventName())
	}
	if note.NotificationType != "shell_completed" || note.Message != "shell done" {
		t.Fatalf("Notification=%+v", note)
	}
}

func TestDecode_Matrix(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		eventHint string
		eventName string
		check     func(t *testing.T, ev copilot.Event)
	}{
		{
			name:      "sessionStart camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","source":"startup","initialPrompt":"hi"}`,
			eventHint: "sessionStart",
			eventName: copilot.EventSessionStart,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.SessionStart)
				if !ok || e.Source != "startup" || e.InitialPrompt() != "hi" {
					t.Fatalf("SessionStart=%+v", ev)
				}
			},
		},
		{
			name:      "sessionStart vscode",
			raw:       `{"hook_event_name":"SessionStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","source":"new","initial_prompt":"go"}`,
			eventName: copilot.EventSessionStart,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.SessionStart)
				if !ok || e.Source != "new" || e.InitialPrompt() != "go" {
					t.Fatalf("SessionStart=%+v", ev)
				}
			},
		},
		{
			name:      "sessionEnd camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","reason":"complete"}`,
			eventHint: "sessionEnd",
			eventName: copilot.EventSessionEnd,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.SessionEnd)
				if !ok || e.Reason != "complete" {
					t.Fatalf("SessionEnd=%+v", ev)
				}
			},
		},
		{
			name:      "userPromptSubmitted camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","prompt":"hello"}`,
			eventHint: "userPromptSubmitted",
			eventName: copilot.EventUserPromptSubmitted,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.UserPromptSubmitted)
				if !ok || e.Prompt != "hello" {
					t.Fatalf("UserPromptSubmitted=%+v", ev)
				}
			},
		},
		{
			name:      "postToolUse camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"view","toolArgs":{},"toolResult":{"resultType":"success","textResultForLlm":"contents"}}`,
			eventHint: "postToolUse",
			eventName: copilot.EventPostToolUse,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.PostToolUse)
				if !ok || e.NativeToolName() != "view" || e.ResultText() != "contents" {
					t.Fatalf("PostToolUse=%+v", ev)
				}
			},
		},
		{
			name:      "postToolUse vscode",
			raw:       `{"hook_event_name":"PostToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"Read","tool_input":{},"tool_result":{"result_type":"success","text_result_for_llm":"file"}}`,
			eventName: copilot.EventPostToolUse,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.PostToolUse)
				if !ok || e.NativeToolName() != "Read" || e.ResultText() != "file" {
					t.Fatalf("PostToolUse=%+v", ev)
				}
			},
		},
		{
			name:      "postToolUseFailure camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"bash","toolArgs":{},"error":"timeout"}`,
			eventHint: "postToolUseFailure",
			eventName: copilot.EventPostToolUseFailure,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.PostToolUseFailure)
				if !ok || e.ErrorMessage() != "timeout" {
					t.Fatalf("PostToolUseFailure=%+v", ev)
				}
			},
		},
		{
			name:      "permissionRequest camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"create","toolArgs":{"path":"a.txt"}}`,
			eventHint: "permissionRequest",
			eventName: copilot.EventPermissionRequest,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.PermissionRequest)
				if !ok || e.NativeToolName() != "create" {
					t.Fatalf("PermissionRequest=%+v", ev)
				}
			},
		},
		{
			name:      "subagentStart camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","transcriptPath":"/t","agentName":"explore","agentDisplayName":"Explore","agentDescription":"search codebase"}`,
			eventHint: "subagentStart",
			eventName: copilot.EventSubagentStart,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.SubagentStart)
				if !ok || e.Name() != "explore" || e.DisplayName() != "Explore" {
					t.Fatalf("SubagentStart=%+v", ev)
				}
			},
		},
		{
			name:      "subagentStop vscode",
			raw:       `{"hook_event_name":"SubagentStop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","stop_reason":"end_turn"}`,
			eventName: copilot.EventSubagentStop,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.SubagentStop)
				if !ok || e.Name() != "task" || e.Reason() != "end_turn" {
					t.Fatalf("SubagentStop=%+v", ev)
				}
			},
		},
		{
			name:      "agentStop camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","transcriptPath":"/t","stopReason":"end_turn"}`,
			eventHint: "agentStop",
			eventName: copilot.EventAgentStop,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.AgentStop)
				if !ok || e.Reason() != "end_turn" {
					t.Fatalf("AgentStop=%+v", ev)
				}
			},
		},
		{
			name:      "preCompact vscode",
			raw:       `{"hook_event_name":"PreCompact","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","trigger":"auto","custom_instructions":"keep tests"}`,
			eventName: copilot.EventPreCompact,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.PreCompact)
				if !ok || e.Trigger != "auto" || e.Instructions() != "keep tests" {
					t.Fatalf("PreCompact=%+v", ev)
				}
			},
		},
		{
			name:      "errorOccurred camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","error":{"message":"slow down","name":"RateLimit"},"errorContext":"model_call","recoverable":true}`,
			eventHint: "errorOccurred",
			eventName: copilot.EventErrorOccurred,
			check: func(t *testing.T, ev copilot.Event) {
				e, ok := ev.(copilot.ErrorOccurred)
				if !ok {
					t.Fatalf("ErrorOccurred=%+v", ev)
				}
				detail, ok := e.Detail()
				if !ok || detail.Name != "RateLimit" || detail.Message != "slow down" {
					t.Fatalf("Detail=%+v ok=%v", detail, ok)
				}
				if e.Recoverable == nil || !*e.Recoverable {
					t.Fatalf("Recoverable=%v", e.Recoverable)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []copilot.Option
			if tt.eventHint != "" {
				opts = append(opts, copilot.WithEvent(tt.eventHint))
			}
			ev, err := copilot.Decode([]byte(tt.raw), opts...)
			if err != nil {
				t.Fatal(err)
			}
			if ev.EventName() != tt.eventName {
				t.Fatalf("EventName=%q, want %q", ev.EventName(), tt.eventName)
			}
			tt.check(t, ev)
		})
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	_, err := copilot.Decode([]byte("not json"), copilot.WithEvent("preToolUse"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "decode payload") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"bash","toolArgs":{"command":"ls"}}`)
	ev, err := copilot.Decode(raw, copilot.WithEvent("preToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(copilot.PreToolUse)
	rawCopy := bytes.Clone(copilot.RawBytes(ev))
	got := pre.ToolArgs.Raw()
	got[0] = 'X'
	if !bytes.Equal(copilot.RawBytes(ev), rawCopy) {
		t.Fatal("mutating ToolArgs.Raw() copy affected RawBytes")
	}
	if bytes.Equal(pre.ToolArgs.Raw(), got) {
		t.Fatal("ToolArgs.Raw() did not return a defensive copy")
	}
}

func TestEncode_PreToolAllowModifiedArgs(t *testing.T) {
	out, code, err := copilot.Encode(copilot.EventPreToolUse, copilot.PreToolResultsForTest().Allow().WithModifiedArgs(map[string]any{"command": "echo safe"}))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"permissionDecision":"allow"`) || !strings.Contains(string(out), `"modifiedArgs"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PostToolUpdatedOutput(t *testing.T) {
	out, code, err := copilot.Encode(copilot.EventPostToolUse, copilot.PostToolResultsForTest().Context("extra guidance").WithModifiedResult("rewritten"))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	s := string(out)
	if !strings.Contains(s, `"textResultForLlm":"rewritten"`) || !strings.Contains(s, `"additionalContext":"extra guidance"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_StopFollowUp(t *testing.T) {
	out, code, err := copilot.Encode(copilot.EventAgentStop, copilot.StopResultsForTest().FollowUp("run the tests"))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestDenyInterrupt(t *testing.T) {
	out, code, err := copilot.Encode(copilot.EventPermissionRequest, copilot.PermissionRequestResultsForTest().Deny("blocked").WithInterrupt(true))
	if err != nil {
		t.Fatal(err)
	}
	if code != copilot.WarnExit {
		t.Fatalf("code=%d, want %d", code, copilot.WarnExit)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) || !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestAsk(t *testing.T) {
	out, code, err := copilot.Encode(copilot.EventPermissionRequest, copilot.PermissionRequestResultsForTest().Ask("needs user confirmation"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PostToolFailureContext(t *testing.T) {
	out, code, err := copilot.Encode(copilot.EventPostToolUseFailure, copilot.PostToolFailureResultsForTest().Context("retry with smaller input"))
	if err != nil {
		t.Fatal(err)
	}
	if code != copilot.WarnExit {
		t.Fatalf("code=%d, want %d", code, copilot.WarnExit)
	}
	if string(out) != "retry with smaller input" {
		t.Fatalf("stdout=%q", out)
	}
}

func TestEncode_SessionStartContext(t *testing.T) {
	out, code, err := copilot.Encode(copilot.EventSessionStart, copilot.SessionStartResultsForTest().Context("project uses go test ./..."))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"additionalContext"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := copilot.Encode(copilot.EventPreToolUse, nil)
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero output should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestEncode_EventOutputMismatch(t *testing.T) {
	_, _, err := copilot.Encode(copilot.EventPostToolUse, copilot.PreToolResultsForTest().Allow())
	if err == nil {
		t.Fatal("expected incompatible event/output error")
	}
}

func TestErrorOccurred_DetailNull(t *testing.T) {
	e := copilot.ErrorOccurred{Error: json.RawMessage("null")}
	if _, ok := e.Detail(); ok {
		t.Fatal("JSON null error payload should be absent")
	}
}

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    int64
		wantErr bool
	}{
		{name: "ms epoch", raw: `1760000000000`, want: 1760000000000},
		{name: "iso8601", raw: `"2026-07-12T10:00:00Z"`, want: 1783850400000},
		{name: "invalid RFC3339", raw: `"not-a-time"`, wantErr: true},
		{name: "non-integer epoch", raw: `1.5`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts copilot.Timestamp
			err := json.Unmarshal([]byte(tt.raw), &ts)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if ts.UnixMilli() != tt.want {
				t.Fatalf("UnixMilli()=%d, want %d", ts.UnixMilli(), tt.want)
			}
		})
	}
}

func TestMux_Serve_PreToolHandlerError(t *testing.T) {
	run.Reset()
	copilot.ResetHandlers()
	new(copilot.Chain).PreToolUse(func(ctx context.Context, hook copilot.Hook[copilot.PreToolUse], _ copilot.PreToolResults) (copilot.PreToolOutput, error) {
		return nil, errors.New("boom")
	})
	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(copilotCamelPreToolUse), &stdout, &bytes.Buffer{}, run.WithEvent("preToolUse"))
	if code != copilot.PreToolErrorExit {
		t.Fatalf("exit = %d, want %d", code, copilot.PreToolErrorExit)
	}
}

func TestMux_Serve_PreToolDeny(t *testing.T) {
	run.Reset()
	copilot.ResetHandlers()
	new(copilot.Chain).PreToolUse(func(ctx context.Context, hook copilot.Hook[copilot.PreToolUse], r copilot.PreToolResults) (copilot.PreToolOutput, error) {
		return r.Deny("nope"), nil
	})
	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(copilotCamelPreToolUse), &stdout, &bytes.Buffer{}, run.WithEvent("preToolUse"))
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stdout.String(), `"permissionDecision":"deny"`) {
		t.Fatalf("stdout=%s", stdout.String())
	}
}

func TestMux_OnDuplicatePanics(t *testing.T) {
	run.Reset()
	copilot.ResetHandlers()
	new(copilot.Chain).PreToolUse(func(ctx context.Context, hook copilot.Hook[copilot.PreToolUse], _ copilot.PreToolResults) (copilot.PreToolOutput, error) {
		return nil, nil
	})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate handler registration")
		}
	}()
	new(copilot.Chain).PreToolUse(func(ctx context.Context, hook copilot.Hook[copilot.PreToolUse], _ copilot.PreToolResults) (copilot.PreToolOutput, error) {
		return nil, nil
	})
}

func TestToolInput_AsBash(t *testing.T) {
	ev, err := copilot.Decode([]byte(`{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"Bash","toolArgs":{"command":"ls -la"}}`), copilot.WithEvent("preToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(copilot.PreToolUse)
	input, ok := pre.Input().AsBash()
	if !ok || input.Command != "ls -la" {
		t.Fatalf("AsBash = %+v, %v", input, ok)
	}
}

func TestRawBytes_PreservesTypedDecode(t *testing.T) {
	ev, err := copilot.Decode([]byte(copilotCamelPreToolUse), copilot.WithEvent("preToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(copilot.RawBytes(ev), []byte(copilotCamelPreToolUse)) {
		t.Fatal("RawBytes mismatch")
	}
}

func TestDecode_VSCodeStopWithAgentScope(t *testing.T) {
	raw := `{
  "hook_event_name": "Stop",
  "session_id": "s2",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "agent_name": "task",
  "stop_reason": "end_turn"
}`
	ev, err := copilot.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(copilot.SubagentStop)
	if !ok {
		t.Fatalf("got %T, want SubagentStop", ev)
	}
	if stop.EventName() != copilot.EventSubagentStop {
		t.Fatalf("EventName=%q", stop.EventName())
	}
	if stop.ReceivedName() != "Stop" {
		t.Fatalf("ReceivedName=%q", stop.ReceivedName())
	}
}

func TestDecode_PostToolUseResultRawPreservesWireShape(t *testing.T) {
	raw := `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"view","toolArgs":{},"toolResult":{"resultType":"success","textResultForLlm":"contents"}}`
	ev, err := copilot.Decode([]byte(raw), copilot.WithEvent("postToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	post := ev.(copilot.PostToolUse)
	got := string(post.ResultRaw())
	if strings.Contains(got, "result_type") || strings.Contains(got, "text_result_for_llm") {
		t.Fatalf("ResultRaw mixed conventions: %s", got)
	}
	if !strings.Contains(got, "textResultForLlm") || !strings.Contains(got, "contents") {
		t.Fatalf("ResultRaw=%s", got)
	}
}

func TestSniffFormat(t *testing.T) {
	if got := copilot.SniffFormat([]byte(copilotCamelPreToolUse)); got != copilot.FormatCamel {
		t.Fatalf("Format=%v, want FormatCamel", got)
	}
	if got := copilot.SniffFormat([]byte(copilotVSCodeStop)); got != copilot.FormatVSCode {
		t.Fatalf("Format=%v, want FormatVSCode", got)
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
