package cursorcore

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
	// fields on HookDataCommon yet; encoding/json ignores them until we model them explicitly.
	raw := []byte(input)
	before := append([]byte(nil), raw...)
	hookData, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	if !bytes.Equal(before, raw) {
		t.Fatalf("NewHookDataCommon must not modify rawJSON bytes")
	}

	assertStringEqual(t, "conv-1", hookData.ConversationID)
	assertStringEqual(t, "gen-1", hookData.GenerationID)
	assertStringEqual(t, "claude-sonnet-4", hookData.Model)
	assertStringEqual(t, "afterFileEdit", hookData.HookEventName)
	assertStringEqual(t, "1.7.2", hookData.CursorVersion)
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

	hookData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	if hookData.UserEmail != nil || hookData.TranscriptPath != nil {
		t.Fatalf("want nil optional pointers, got user=%#v path=%#v", hookData.UserEmail, hookData.TranscriptPath)
	}
}

func TestNewHookDataCommon_hookEventName(t *testing.T) {
	hookData, err := NewHookDataCommon([]byte(`{"hook_event_name":"afterFileEdit"}`))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	assertStringEqual(t, "afterFileEdit", hookData.HookEventName)
}

func TestNewHookDataCommon_malformedJSON(t *testing.T) {
	raw := []byte(`not json`)
	before := append([]byte(nil), raw...)
	hookData, err := NewHookDataCommon(raw)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
	if !bytes.Equal(before, raw) {
		t.Fatalf("NewHookDataCommon must not modify rawJSON bytes")
	}
	assertHookDataCommonZero(t, hookData)
}

func TestNewHookDataCommon_emptyInput(t *testing.T) {
	raw := []byte{}
	before := append([]byte(nil), raw...)
	hookData, err := NewHookDataCommon(raw)
	if err == nil {
		t.Fatal("expected error for empty JSON input")
	}
	if !bytes.Equal(before, raw) {
		t.Fatalf("NewHookDataCommon must not modify rawJSON bytes")
	}
	assertHookDataCommonZero(t, hookData)
}

func TestNewHookDataCommon_emptyObject(t *testing.T) {
	// Valid object with no conversation_id, hook_event_name, or other known keys — json leaves zero values.
	input := `{}`
	raw := []byte(input)
	before := append([]byte(nil), raw...)
	hookData, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	if !bytes.Equal(before, raw) {
		t.Fatalf("NewHookDataCommon must not modify rawJSON bytes")
	}
	assertStringEqual(t, "", hookData.HookEventName)
	assertStringEqual(t, "", hookData.ConversationID)
	assertStringEqual(t, "", hookData.GenerationID)
	assertStringEqual(t, "", hookData.Model)
	assertStringEqual(t, "", hookData.CursorVersion)
	if hookData.WorkspaceRoots != nil {
		t.Fatalf("WorkspaceRoots: want nil, got %#v", hookData.WorkspaceRoots)
	}
	if hookData.UserEmail != nil || hookData.TranscriptPath != nil {
		t.Fatalf("optional pointers: want nil, got user=%#v path=%#v", hookData.UserEmail, hookData.TranscriptPath)
	}
}

func assertStringEqual(t *testing.T, want, got string) {
	t.Helper()
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func assertHookDataCommonZero(t *testing.T, hookData HookDataCommon) {
	t.Helper()
	want := HookDataCommon{}
	if !reflect.DeepEqual(want, hookData) {
		t.Fatalf("want zero HookDataCommon, got %#v", hookData)
	}
}
