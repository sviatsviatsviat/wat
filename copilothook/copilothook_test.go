package copilothook_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/copilothook"
	"github.com/sviatsviatsviat/wat/copilothook/tools"
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
	ev, err := copilothook.Decode([]byte(copilotCamelPreToolUse), copilothook.WithEvent("preToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(copilothook.PreToolUse)
	if !ok || pre.Session() != "s1" || pre.Cwd != "/w" {
		t.Fatalf("bad event: %+v", ev)
	}
	if pre.NativeToolName() != "bash" || pre.ShellCommand() != "rm -rf /" {
		t.Fatalf("bad tool: name=%q shell=%q", pre.NativeToolName(), pre.ShellCommand())
	}
	if !bytes.Equal(copilothook.RawBytes(ev), []byte(copilotCamelPreToolUse)) {
		t.Fatal("Raw not preserved")
	}

	out, code, err := copilothook.Encode(copilothook.EventPreToolUse, copilothook.PreToolOutput{
		Decision: copilothook.DecisionDeny,
		Reason:   "destructive command",
	})
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
	_, err := copilothook.Decode([]byte(copilotCamelPreToolUse))
	if err == nil {
		t.Fatal("expected error without WithEvent")
	}
	if !strings.Contains(err.Error(), "WithEvent") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_VSCodeStop(t *testing.T) {
	ev, err := copilothook.Decode([]byte(copilotVSCodeStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(copilothook.AgentStop)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if stop.EventName() != copilothook.EventAgentStop {
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
	ev, err := copilothook.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(copilothook.PreToolUse)
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
	ev, err := copilothook.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	note, ok := ev.(copilothook.Notification)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if note.EventName() != copilothook.EventNotification {
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
		check     func(t *testing.T, ev copilothook.Event)
	}{
		{
			name:      "sessionStart camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","source":"startup","initialPrompt":"hi"}`,
			eventHint: "sessionStart",
			eventName: copilothook.EventSessionStart,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.SessionStart)
				if !ok || e.Source != "startup" || e.InitialPrompt() != "hi" {
					t.Fatalf("SessionStart=%+v", ev)
				}
			},
		},
		{
			name:      "sessionStart vscode",
			raw:       `{"hook_event_name":"SessionStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","source":"new","initial_prompt":"go"}`,
			eventName: copilothook.EventSessionStart,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.SessionStart)
				if !ok || e.Source != "new" || e.InitialPrompt() != "go" {
					t.Fatalf("SessionStart=%+v", ev)
				}
			},
		},
		{
			name:      "sessionEnd camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","reason":"complete"}`,
			eventHint: "sessionEnd",
			eventName: copilothook.EventSessionEnd,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.SessionEnd)
				if !ok || e.Reason != "complete" {
					t.Fatalf("SessionEnd=%+v", ev)
				}
			},
		},
		{
			name:      "userPromptSubmitted camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","prompt":"hello"}`,
			eventHint: "userPromptSubmitted",
			eventName: copilothook.EventUserPromptSubmitted,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.UserPromptSubmitted)
				if !ok || e.Prompt != "hello" {
					t.Fatalf("UserPromptSubmitted=%+v", ev)
				}
			},
		},
		{
			name:      "postToolUse camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"view","toolArgs":{},"toolResult":{"resultType":"success","textResultForLlm":"contents"}}`,
			eventHint: "postToolUse",
			eventName: copilothook.EventPostToolUse,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.PostToolUse)
				if !ok || e.NativeToolName() != "view" || e.ResultText() != "contents" {
					t.Fatalf("PostToolUse=%+v", ev)
				}
			},
		},
		{
			name:      "postToolUse vscode",
			raw:       `{"hook_event_name":"PostToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"Read","tool_input":{},"tool_result":{"result_type":"success","text_result_for_llm":"file"}}`,
			eventName: copilothook.EventPostToolUse,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.PostToolUse)
				if !ok || e.NativeToolName() != "Read" || e.ResultText() != "file" {
					t.Fatalf("PostToolUse=%+v", ev)
				}
			},
		},
		{
			name:      "postToolUseFailure camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"bash","toolArgs":{},"error":"timeout"}`,
			eventHint: "postToolUseFailure",
			eventName: copilothook.EventPostToolUseFailure,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.PostToolUseFailure)
				if !ok || e.ErrorMessage() != "timeout" {
					t.Fatalf("PostToolUseFailure=%+v", ev)
				}
			},
		},
		{
			name:      "permissionRequest camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"create","toolArgs":{"path":"a.txt"}}`,
			eventHint: "permissionRequest",
			eventName: copilothook.EventPermissionRequest,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.PermissionRequest)
				if !ok || e.NativeToolName() != "create" {
					t.Fatalf("PermissionRequest=%+v", ev)
				}
			},
		},
		{
			name:      "subagentStart camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","transcriptPath":"/t","agentName":"explore","agentDisplayName":"Explore","agentDescription":"search codebase"}`,
			eventHint: "subagentStart",
			eventName: copilothook.EventSubagentStart,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.SubagentStart)
				if !ok || e.Name() != "explore" || e.DisplayName() != "Explore" {
					t.Fatalf("SubagentStart=%+v", ev)
				}
			},
		},
		{
			name:      "subagentStop vscode",
			raw:       `{"hook_event_name":"SubagentStop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","stop_reason":"end_turn"}`,
			eventName: copilothook.EventSubagentStop,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.SubagentStop)
				if !ok || e.Name() != "task" || e.Reason() != "end_turn" {
					t.Fatalf("SubagentStop=%+v", ev)
				}
			},
		},
		{
			name:      "agentStop camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","transcriptPath":"/t","stopReason":"end_turn"}`,
			eventHint: "agentStop",
			eventName: copilothook.EventAgentStop,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.AgentStop)
				if !ok || e.Reason() != "end_turn" {
					t.Fatalf("AgentStop=%+v", ev)
				}
			},
		},
		{
			name:      "preCompact vscode",
			raw:       `{"hook_event_name":"PreCompact","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","trigger":"auto","custom_instructions":"keep tests"}`,
			eventName: copilothook.EventPreCompact,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.PreCompact)
				if !ok || e.Trigger != "auto" || e.Instructions() != "keep tests" {
					t.Fatalf("PreCompact=%+v", ev)
				}
			},
		},
		{
			name:      "errorOccurred camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","error":{"message":"slow down","name":"RateLimit"},"errorContext":"model_call","recoverable":true}`,
			eventHint: "errorOccurred",
			eventName: copilothook.EventErrorOccurred,
			check: func(t *testing.T, ev copilothook.Event) {
				e, ok := ev.(copilothook.ErrorOccurred)
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
			var opts []copilothook.Option
			if tt.eventHint != "" {
				opts = append(opts, copilothook.WithEvent(tt.eventHint))
			}
			ev, err := copilothook.Decode([]byte(tt.raw), opts...)
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
	_, err := copilothook.Decode([]byte("not json"), copilothook.WithEvent("preToolUse"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "decode payload") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"bash","toolArgs":{"command":"ls"}}`)
	ev, err := copilothook.Decode(raw, copilothook.WithEvent("preToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(copilothook.PreToolUse)
	rawCopy := bytes.Clone(copilothook.RawBytes(ev))
	pre.ToolArgs[0] = '{'
	if !bytes.Equal(copilothook.RawBytes(ev), rawCopy) {
		t.Fatal("mutating ToolArgs affected RawBytes")
	}
}

func TestEncode_PreToolAllowModifiedArgs(t *testing.T) {
	out, code, err := copilothook.Encode(copilothook.EventPreToolUse, copilothook.PreToolOutput{
		Decision:     copilothook.DecisionAllow,
		ModifiedArgs: map[string]any{"command": "echo safe"},
	})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"permissionDecision":"allow"`) || !strings.Contains(string(out), `"modifiedArgs"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PostToolUpdatedOutput(t *testing.T) {
	out, code, err := copilothook.Encode(copilothook.EventPostToolUse, copilothook.PostToolOutput{
		ModifiedResult:    "rewritten",
		AdditionalContext: "extra guidance",
	})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	s := string(out)
	if !strings.Contains(s, `"textResultForLlm":"rewritten"`) || !strings.Contains(s, `"additionalContext":"extra guidance"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_StopFollowUp(t *testing.T) {
	out, code, err := copilothook.Encode(copilothook.EventAgentStop, copilothook.StopOutput{Reason: "run the tests"})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestDenyInterrupt(t *testing.T) {
	out, code, err := copilothook.Encode(copilothook.EventPermissionRequest, copilothook.PermissionRequestOutput{
		Behavior:  "deny",
		Message:   "blocked",
		Interrupt: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if code != copilothook.WarnExit {
		t.Fatalf("code=%d, want %d", code, copilothook.WarnExit)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) || !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestAsk(t *testing.T) {
	out, code, err := copilothook.Encode(copilothook.EventPermissionRequest, copilothook.PermissionRequestOutput{
		Behavior:         "deny",
		Message:          "needs user confirmation",
		SuppressWarnExit: true,
	})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PostToolFailureContext(t *testing.T) {
	out, code, err := copilothook.Encode(copilothook.EventPostToolUseFailure, copilothook.PostToolFailureOutput{
		Context: "retry with smaller input",
	})
	if err != nil {
		t.Fatal(err)
	}
	if code != copilothook.WarnExit {
		t.Fatalf("code=%d, want %d", code, copilothook.WarnExit)
	}
	if string(out) != "retry with smaller input" {
		t.Fatalf("stdout=%q", out)
	}
}

func TestEncode_SessionStartContext(t *testing.T) {
	out, code, err := copilothook.Encode(copilothook.EventSessionStart, copilothook.SessionStartOutput{
		AdditionalContext: "project uses go test ./...",
	})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"additionalContext"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := copilothook.Encode(copilothook.EventPreToolUse, copilothook.PreToolOutput{})
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero output should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestEncode_EventOutputMismatch(t *testing.T) {
	_, _, err := copilothook.Encode(copilothook.EventPostToolUse, copilothook.PreToolOutput{
		Decision: copilothook.DecisionAllow,
	})
	if err == nil {
		t.Fatal("expected incompatible event/output error")
	}
}

func TestErrorOccurred_DetailNull(t *testing.T) {
	e := copilothook.ErrorOccurred{Error: json.RawMessage("null")}
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
			var ts copilothook.Timestamp
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
	mux := copilothook.NewMux()
	copilothook.On(mux, func(ctx context.Context, ev copilothook.PreToolUse) (copilothook.PreToolOutput, error) {
		return copilothook.PreToolOutput{}, errors.New("boom")
	})
	var stdout bytes.Buffer
	code := mux.Serve(context.Background(), strings.NewReader(copilotCamelPreToolUse), &stdout, &bytes.Buffer{}, copilothook.WithEvent("preToolUse"))
	if code != copilothook.PreToolErrorExit {
		t.Fatalf("exit = %d, want %d", code, copilothook.PreToolErrorExit)
	}
}

func TestMux_Serve_PreToolDeny(t *testing.T) {
	mux := copilothook.NewMux()
	copilothook.On(mux, func(ctx context.Context, ev copilothook.PreToolUse) (copilothook.PreToolOutput, error) {
		return copilothook.PreToolOutput{Decision: copilothook.DecisionDeny, Reason: "nope"}, nil
	})
	var stdout bytes.Buffer
	code := mux.Serve(context.Background(), strings.NewReader(copilotCamelPreToolUse), &stdout, &bytes.Buffer{}, copilothook.WithEvent("preToolUse"))
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stdout.String(), `"permissionDecision":"deny"`) {
		t.Fatalf("stdout=%s", stdout.String())
	}
}

func TestToolInputAs_Bash(t *testing.T) {
	input, err := tools.ToolInputAs[tools.BashInput](json.RawMessage(`{"command":"ls -la"}`))
	if err != nil {
		t.Fatal(err)
	}
	if input.Command != "ls -la" {
		t.Fatalf("Command=%q", input.Command)
	}
}

func TestRawBytes_PreservesTypedDecode(t *testing.T) {
	ev, err := copilothook.Decode([]byte(copilotCamelPreToolUse), copilothook.WithEvent("preToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(copilothook.RawBytes(ev), []byte(copilotCamelPreToolUse)) {
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
	ev, err := copilothook.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(copilothook.SubagentStop)
	if !ok {
		t.Fatalf("got %T, want SubagentStop", ev)
	}
	if stop.EventName() != copilothook.EventSubagentStop {
		t.Fatalf("EventName=%q", stop.EventName())
	}
	if stop.ReceivedName() != "Stop" {
		t.Fatalf("ReceivedName=%q", stop.ReceivedName())
	}
}

func TestDecode_PostToolUseResultRawPreservesWireShape(t *testing.T) {
	raw := `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"view","toolArgs":{},"toolResult":{"resultType":"success","textResultForLlm":"contents"}}`
	ev, err := copilothook.Decode([]byte(raw), copilothook.WithEvent("postToolUse"))
	if err != nil {
		t.Fatal(err)
	}
	post := ev.(copilothook.PostToolUse)
	got := string(post.ResultRaw())
	if strings.Contains(got, "result_type") || strings.Contains(got, "text_result_for_llm") {
		t.Fatalf("ResultRaw mixed conventions: %s", got)
	}
	if !strings.Contains(got, "textResultForLlm") || !strings.Contains(got, "contents") {
		t.Fatalf("ResultRaw=%s", got)
	}
}

func TestSniffFormat(t *testing.T) {
	if got := copilothook.SniffFormat([]byte(copilotCamelPreToolUse)); got != copilothook.FormatCamel {
		t.Fatalf("Format=%v, want FormatCamel", got)
	}
	if got := copilothook.SniffFormat([]byte(copilotVSCodeStop)); got != copilothook.FormatVSCode {
		t.Fatalf("Format=%v, want FormatVSCode", got)
	}
}
