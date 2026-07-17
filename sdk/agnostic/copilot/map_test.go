package copilot

import (
	"bytes"
	"encoding/json"
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

func mapRaw(t *testing.T, raw []byte, eventHint string) *model.Event {
	t.Helper()
	var opts []sdkcopilot.Option
	if eventHint != "" {
		opts = append(opts, sdkcopilot.WithEvent(eventHint))
	}
	native, err := sdkcopilot.Decode(raw, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return MapEvent(native, raw)
}

func TestCopilotMapEvent_VSCodeStop(t *testing.T) {
	ev := mapRaw(t, []byte(copilotVSCodeStop), "")
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

func TestCopilotMapEvent_VSCodePreToolBash(t *testing.T) {
	raw := `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "Bash",
  "tool_input": {"command": "ls -la"}
}`
	ev := mapRaw(t, []byte(raw), "")
	if ev.Tool == nil || ev.Tool.Name != model.ToolBash || ev.Tool.Native != "Bash" || ev.Tool.Shell != "ls -la" {
		t.Fatalf("Tool=%+v", ev.Tool)
	}
}

func TestCopilotMapEvent_NotificationWithoutEventHint(t *testing.T) {
	raw := `{
  "sessionId": "s1",
  "timestamp": 1760000000000,
  "cwd": "/w",
  "hook_event_name": "Notification",
  "message": "shell done",
  "title": "Shell completed",
  "notification_type": "shell_completed"
}`
	ev := mapRaw(t, []byte(raw), "")
	if ev.Kind != model.KindNotification || ev.Name != "Notification" {
		t.Fatalf("model.Kind=%v Name=%q", ev.Kind, ev.Name)
	}
	if ev.Note == nil || ev.Note.Type != "shell_completed" || ev.Note.Message != "shell done" {
		t.Fatalf("Note=%+v", ev.Note)
	}
}

func TestCopilotMapEvent_Matrix(t *testing.T) {
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
				if ev.Subagent == nil || ev.Subagent.Type != "explore" || ev.Subagent.Task != "Explore" || ev.Subagent.Summary != "search codebase" {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
			},
		},
		{
			name: "subagentStop vscode",
			raw:  `{"hook_event_name":"SubagentStop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","stop_reason":"end_turn"}`,
			kind: model.KindSubagentStop,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Subagent == nil || ev.Subagent.Type != "task" || ev.Subagent.Status != "end_turn" {
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
		{
			name:      "preToolUse camel",
			raw:       copilotCamelPreToolUse,
			eventHint: "preToolUse",
			kind:      model.KindPreTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolBash || ev.Tool.Shell != "rm -rf /" {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := mapRaw(t, []byte(tt.raw), tt.eventHint)
			if ev.Kind != tt.kind {
				t.Fatalf("model.Kind=%v, want %v", ev.Kind, tt.kind)
			}
			tt.check(t, ev)
		})
	}
}

func TestCopilotMapEvent_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"sessionId":"s","timestamp":1,"cwd":"/w","toolName":"bash","toolArgs":{"command":"ls"}}`)
	ev := mapRaw(t, raw, "preToolUse")
	rawCopy := bytes.Clone(ev.Raw)
	got := ev.Tool.Input.Raw()
	got[0] = 'X'
	if !bytes.Equal(ev.Raw, rawCopy) {
		t.Fatal("mutating Tool.Input.Raw() copy affected Event.Raw")
	}
	if bytes.Equal(ev.Tool.Input.Raw(), got) {
		t.Fatal("Tool.Input.Raw() did not return a defensive copy")
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
