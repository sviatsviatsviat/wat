package cursor

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewHookDataCommon_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "afterFileEdit",
		"conversation_id": "conv-1",
		"generation_id": "gen-1",
		"model": "claude-sonnet-4",
		"cursor_version": "1.7.2",
		"workspace_roots": ["/repo", "/repo-2"],
		"user_email": "dev@example.com",
		"transcript_path": "/tmp/transcript.jsonl",
		"session_id": "sess-hook",
		"cwd": "/repo",
		"tool_name": "Shell",
		"command": "npm install",
		"path": "src/main.ts"
	}`

	// session_id, cwd, tool_name, command, and path appear in some Cursor payloads but are not
	// fields on hookDataCommon yet; encoding/json ignores them until we model them explicitly.
	raw := []byte(input)
	hookData, err := newHookDataCommon(raw)
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	if string(raw) != input {
		t.Fatalf("newHookDataCommon must not replace rawJSON bytes")
	}

	assertEqual(t, "conv-1", hookData.ConversationID)
	assertEqual(t, "gen-1", hookData.GenerationID)
	assertEqual(t, "claude-sonnet-4", hookData.Model)
	assertEqual(t, "afterFileEdit", hookData.HookEventName)
	assertEqual(t, "1.7.2", hookData.CursorVersion)
	wantRoots := []string{"/repo", "/repo-2"}
	if !reflect.DeepEqual(wantRoots, hookData.WorkspaceRoots) {
		t.Fatalf("WorkspaceRoots: want %#v, got %#v", wantRoots, hookData.WorkspaceRoots)
	}
	if hookData.UserEmail == nil || *hookData.UserEmail != "dev@example.com" {
		t.Fatalf("UserEmail: want pointer to dev@example.com, got %#v", hookData.UserEmail)
	}
	if hookData.TranscriptPath == nil || *hookData.TranscriptPath != "/tmp/transcript.jsonl" {
		t.Fatalf("TranscriptPath: want pointer to /tmp/transcript.jsonl, got %#v", hookData.TranscriptPath)
	}
}

func TestNewHookDataCommon_nullOptionalStrings(t *testing.T) {
	input := `{
		"hook_event_name": "sessionEnd",
		"conversation_id": "c",
		"user_email": null,
		"transcript_path": null
	}`

	hookData, err := newHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	if hookData.UserEmail != nil || hookData.TranscriptPath != nil {
		t.Fatalf("want nil optional pointers, got user=%#v path=%#v", hookData.UserEmail, hookData.TranscriptPath)
	}
}

func TestNewHookDataCommon_hookEventName(t *testing.T) {
	hookData, err := newHookDataCommon([]byte(`{"hook_event_name":"afterFileEdit"}`))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	assertEqual(t, "afterFileEdit", hookData.HookEventName)
}

func TestNewHookDataCommon_malformedJSON(t *testing.T) {
	raw := []byte(`not json`)
	before := append([]byte(nil), raw...)
	hookData, err := newHookDataCommon(raw)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
	if !bytes.Equal(before, raw) {
		t.Fatalf("newHookDataCommon must not modify rawJSON bytes")
	}
	assertHookDataCommonZero(t, hookData)
}

func TestNewHookDataCommon_emptyInput(t *testing.T) {
	raw := []byte{}
	before := append([]byte(nil), raw...)
	hookData, err := newHookDataCommon(raw)
	if err == nil {
		t.Fatal("expected error for empty JSON input")
	}
	if !bytes.Equal(before, raw) {
		t.Fatalf("newHookDataCommon must not modify rawJSON bytes")
	}
	assertHookDataCommonZero(t, hookData)
}

func TestNewHookDataCommon_emptyObject(t *testing.T) {
	// Valid object with no conversation_id, hook_event_name, or other known keys — json leaves zero values.
	input := `{}`
	raw := []byte(input)
	hookData, err := newHookDataCommon(raw)
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	if string(raw) != input {
		t.Fatalf("newHookDataCommon must not replace rawJSON bytes")
	}
	assertEqual(t, "", hookData.HookEventName)
	assertEqual(t, "", hookData.ConversationID)
	assertEqual(t, "", hookData.GenerationID)
	assertEqual(t, "", hookData.Model)
	assertEqual(t, "", hookData.CursorVersion)
	if hookData.WorkspaceRoots != nil {
		t.Fatalf("WorkspaceRoots: want nil, got %#v", hookData.WorkspaceRoots)
	}
	if hookData.UserEmail != nil || hookData.TranscriptPath != nil {
		t.Fatalf("optional pointers: want nil, got user=%#v path=%#v", hookData.UserEmail, hookData.TranscriptPath)
	}
}

func assertEqual(t *testing.T, want, got string) {
	t.Helper()
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

// assertHookDataCommonZero checks the hookDataCommon zero value returned on Unmarshal error.
func assertHookDataCommonZero(t *testing.T, hookData hookDataCommon) {
	t.Helper()
	want := hookDataCommon{}
	if !reflect.DeepEqual(want, hookData) {
		t.Fatalf("want zero hookDataCommon, got %#v", hookData)
	}
}
