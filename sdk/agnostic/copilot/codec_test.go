package copilot

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
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

func TestCopilotDecodeEncode_PreToolDeny(t *testing.T) {
	c := &Codec{}
	ev, err := c.Decode([]byte(copilotCamelPreToolUse), "preToolUse")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != model.KindPreTool || ev.Session != "s1" || ev.Cwd != "/w" {
		t.Fatalf("bad event: %+v", ev)
	}
	if ev.Tool == nil || ev.Tool.Name != model.ToolBash || ev.Tool.Native != "bash" || ev.Tool.Shell != "rm -rf /" {
		t.Fatalf("bad tool: %+v", ev.Tool)
	}
	if !bytes.Equal(ev.Raw, []byte(copilotCamelPreToolUse)) {
		t.Fatal("Raw not preserved")
	}

	out, code, err := c.Encode(ev, model.Deny("destructive command"))
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

func TestCopilotDecode_CamelCaseRequiresEventHint(t *testing.T) {
	c := &Codec{}
	_, err := c.Decode([]byte(copilotCamelPreToolUse), "")
	if err == nil {
		t.Fatal("expected error without eventHint")
	}
	if !errors.Is(err, sdkcopilot.ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
	if !strings.Contains(err.Error(), "eventHint") {
		t.Fatalf("error = %v", err)
	}
}

func TestCopilotDecode_VSCodeStop(t *testing.T) {
	c := &Codec{}
	ev, err := c.Decode([]byte(copilotVSCodeStop), "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != model.KindStop || ev.Name != "Stop" {
		t.Fatalf("model.Kind=%v Name=%q", ev.Kind, ev.Name)
	}
	if ev.Turn == nil || ev.Turn.Status != "end_turn" {
		t.Fatalf("Turn=%+v", ev.Turn)
	}
	if ev.TranscriptPath != "/tmp/t" {
		t.Fatalf("TranscriptPath=%q", ev.TranscriptPath)
	}
}

func TestCopilotDecode_VSCodePreToolBash(t *testing.T) {
	raw := `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "Bash",
  "tool_input": {"command": "ls -la"}
}`
	c := &Codec{}
	ev, err := c.Decode([]byte(raw), "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Tool == nil || ev.Tool.Name != model.ToolBash || ev.Tool.Native != "Bash" || ev.Tool.Shell != "ls -la" {
		t.Fatalf("Tool=%+v", ev.Tool)
	}
}

func TestCopilotDecode_NotificationWithoutEventHint(t *testing.T) {
	raw := `{
  "sessionId": "s1",
  "timestamp": 1760000000000,
  "cwd": "/w",
  "hook_event_name": "Notification",
  "message": "shell done",
  "title": "Shell completed",
  "notification_type": "shell_completed"
}`
	c := &Codec{}
	ev, err := c.Decode([]byte(raw), "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != model.KindNotification || ev.Name != "Notification" {
		t.Fatalf("model.Kind=%v Name=%q", ev.Kind, ev.Name)
	}
	if ev.Note == nil || ev.Note.Type != "shell_completed" || ev.Note.Message != "shell done" {
		t.Fatalf("Note=%+v", ev.Note)
	}
}

func TestCopilotDecode_Matrix(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		eventHint string
		kind      model.Kind
		check     func(t *testing.T, ev *model.Event)
	}{
		{
			name:      "sessionStart camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","source":"startup","initialPrompt":"hi"}`,
			eventHint: "sessionStart",
			kind:      model.KindSessionStart,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Life == nil || ev.Life.Source != "startup" || ev.Life.InitialPrompt != "hi" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "sessionStart vscode",
			raw:  `{"hook_event_name":"SessionStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","source":"new","initial_prompt":"go"}`,
			kind: model.KindSessionStart,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Life == nil || ev.Life.Source != "new" || ev.Life.InitialPrompt != "go" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name:      "sessionEnd camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","reason":"complete"}`,
			eventHint: "sessionEnd",
			kind:      model.KindSessionEnd,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Life == nil || ev.Life.Reason != "complete" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name:      "userPromptSubmitted camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","prompt":"hello"}`,
			eventHint: "userPromptSubmitted",
			kind:      model.KindUserPrompt,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Prompt != "hello" {
					t.Fatalf("Prompt=%q", ev.Prompt)
				}
			},
		},
		{
			name:      "postToolUse camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"view","toolArgs":{},"toolResult":{"resultType":"success","textResultForLlm":"contents"}}`,
			eventHint: "postToolUse",
			kind:      model.KindPostTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolRead {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
				if ev.Result == nil || ev.Result.Text != "contents" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "postToolUse vscode",
			raw:  `{"hook_event_name":"PostToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"Read","tool_input":{},"tool_result":{"result_type":"success","text_result_for_llm":"file"}}`,
			kind: model.KindPostTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolRead {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
				if ev.Result == nil || ev.Result.Text != "file" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name:      "postToolUseFailure camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"bash","toolArgs":{},"error":"timeout"}`,
			eventHint: "postToolUseFailure",
			kind:      model.KindPostToolFailure,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Result == nil || ev.Result.Error != "timeout" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name:      "permissionRequest camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"create","toolArgs":{"path":"a.txt"}}`,
			eventHint: "permissionRequest",
			kind:      model.KindPermissionRequest,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolWrite {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name:      "subagentStart camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","transcriptPath":"/t","agentName":"explore","agentDisplayName":"Explore","agentDescription":"search codebase"}`,
			eventHint: "subagentStart",
			kind:      model.KindSubagentStart,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Subagent == nil || ev.Subagent.Type != "explore" || ev.Subagent.Task != "Explore" {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
			},
		},
		{
			name: "subagentStop vscode",
			raw:  `{"hook_event_name":"SubagentStop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","stop_reason":"end_turn"}`,
			kind: model.KindSubagentStop,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Subagent == nil || ev.Subagent.Type != "task" {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
				if ev.Turn == nil || ev.Turn.Status != "end_turn" {
					t.Fatalf("Turn=%+v", ev.Turn)
				}
			},
		},
		{
			name:      "agentStop camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","transcriptPath":"/t","stopReason":"end_turn"}`,
			eventHint: "agentStop",
			kind:      model.KindStop,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Turn == nil || ev.Turn.Status != "end_turn" {
					t.Fatalf("Turn=%+v", ev.Turn)
				}
			},
		},
		{
			name: "preCompact vscode",
			raw:  `{"hook_event_name":"PreCompact","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","trigger":"auto","custom_instructions":"keep tests"}`,
			kind: model.KindPreCompact,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Compact == nil || ev.Compact.Trigger != "auto" || ev.Compact.CustomInstructions != "keep tests" {
					t.Fatalf("Compact=%+v", ev.Compact)
				}
			},
		},
		{
			name:      "errorOccurred camel",
			raw:       `{"sessionId":"s","timestamp":1,"cwd":"/w","error":{"message":"slow down","name":"RateLimit"},"errorContext":"model_call","recoverable":true}`,
			eventHint: "errorOccurred",
			kind:      model.KindAgentError,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Note == nil || ev.Note.Type != "RateLimit" || ev.Note.Message != "slow down" {
					t.Fatalf("Note=%+v", ev.Note)
				}
				if ev.Note.Recoverable == nil || !*ev.Note.Recoverable {
					t.Fatalf("Recoverable=%v", ev.Note.Recoverable)
				}
			},
		},
	}
	c := &Codec{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := c.Decode([]byte(tt.raw), tt.eventHint)
			if err != nil {
				t.Fatal(err)
			}
			if ev.Kind != tt.kind {
				t.Fatalf("model.Kind=%v, want %v", ev.Kind, tt.kind)
			}
			tt.check(t, ev)
		})
	}
}

func TestCopilotDecode_InvalidJSON(t *testing.T) {
	c := &Codec{}
	_, err := c.Decode([]byte("not json"), "preToolUse")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sdkcopilot.ErrUnrecognizedFormat) {
		t.Fatalf("errors.Is ErrUnrecognizedFormat = false, err = %v", err)
	}
	if !strings.Contains(err.Error(), "copilot: decode payload") {
		t.Fatalf("error = %v", err)
	}
}

func TestCopilotDecode_InvalidPayloadPreservesSentinel(t *testing.T) {
	c := &Codec{}
	raw := []byte(`{"sessionId":"s1","timestamp":1,"cwd":[]}`)
	_, err := c.Decode(raw, "preToolUse")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sdkcopilot.ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
	if !strings.HasPrefix(err.Error(), "copilot: decode payload") {
		t.Fatalf("error = %v", err)
	}
}

func TestCopilotDecode_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"bash","toolArgs":{"command":"ls"}}`)
	c := &Codec{}
	ev, err := c.Decode(raw, "preToolUse")
	if err != nil {
		t.Fatal(err)
	}
	rawCopy := bytes.Clone(ev.Raw)
	ev.Tool.Input[0] = '{'
	if !bytes.Equal(ev.Raw, rawCopy) {
		t.Fatal("mutating Tool.Input affected Event.Raw")
	}
}

func TestCopilotEncode_PreToolAllowModifiedArgs(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindPreTool, Name: "preToolUse"}
	out, code, err := c.Encode(ev, model.Result{
		Decision:     model.DecisionAllow,
		UpdatedInput: map[string]any{"command": "echo safe"},
	})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"permissionDecision":"allow"`) || !strings.Contains(string(out), `"modifiedArgs"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCopilotEncode_EmptyNameUsesKind(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindPreTool, Name: ""}
	out, code, err := c.Encode(ev, model.Result{Decision: model.DecisionAllow})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"permissionDecision":"allow"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCopilotEncode_PostToolUpdatedOutput(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindPostTool, Name: "postToolUse"}
	text := "rewritten"
	out, code, err := c.Encode(ev, model.Result{
		UpdatedOutput: &text,
		Context:       "extra guidance",
	})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	s := string(out)
	if !strings.Contains(s, `"textResultForLlm":"rewritten"`) || !strings.Contains(s, `"additionalContext":"extra guidance"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCopilotEncode_StopFollowUp(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindStop, Name: "agentStop"}
	out, code, err := c.Encode(ev, model.Result{FollowUp: "run the tests"})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCopilotEncode_SubagentStopFollowUp(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindSubagentStop, Name: "subagentStop"}
	out, code, err := c.Encode(ev, model.Result{FollowUp: "finish review"})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCopilotEncode_PermissionRequestDenyInterrupt(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindPermissionRequest, Name: "permissionRequest"}
	out, code, err := c.Encode(ev, model.Result{Decision: model.DecisionDeny, Reason: "blocked", HaltSession: true})
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

func TestCopilotEncode_PermissionRequestAsk(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindPermissionRequest, Name: "permissionRequest"}
	out, code, err := c.Encode(ev, model.Ask("needs user confirmation"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCopilotEncode_PostToolFailureContext(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindPostToolFailure, Name: "postToolUseFailure"}
	out, code, err := c.Encode(ev, model.Context("retry with smaller input"))
	if err != nil {
		t.Fatal(err)
	}
	if code != WarnExit {
		t.Fatalf("code=%d, want %d", code, WarnExit)
	}
	if string(out) != "retry with smaller input" {
		t.Fatalf("stdout=%q", out)
	}
}

func TestCopilotEncode_SessionStartContext(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindSessionStart, Name: "sessionStart"}
	out, code, err := c.Encode(ev, model.Context("project uses go test ./..."))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"additionalContext"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCopilotEncode_ZeroResult(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Copilot, Kind: model.KindPreTool, Name: "preToolUse"}
	out, code, err := c.Encode(ev, model.Result{})
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero result should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestCopilotEncode_NilEvent(t *testing.T) {
	c := &Codec{}
	_, _, err := c.Encode(nil, model.Deny("nope"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "nil event") {
		t.Fatalf("error = %v", err)
	}
}

func TestCopilotTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    int64 // unix ms
		wantErr bool
	}{
		{name: "ms epoch", raw: `1760000000000`, want: 1760000000000},
		{name: "iso8601", raw: `"2026-07-12T10:00:00Z"`, want: 1783850400000},
		{name: "invalid RFC3339", raw: `"not-a-time"`, wantErr: true},
		{name: "non-integer epoch", raw: `1.5`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts sdkcopilot.Timestamp
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
