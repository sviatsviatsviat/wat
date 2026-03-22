package cli

import (
	"strings"
	"testing"
)

func TestReadHookStdinJSON_EmptyInput(t *testing.T) {
	hookStdinJSON, err := ReadHookStdinJSON(strings.NewReader("   "))
	if err != nil {
		t.Fatalf("ReadHookStdinJSON() error = %v", err)
	}
	if hookStdinJSON != nil {
		t.Fatalf("expected nil body for whitespace-only stdin, got %q", hookStdinJSON)
	}
}

func TestReadHookStdinJSON_ValidJSON(t *testing.T) {
	input := `{"hook_event_name":"afterFileEdit","conversation_id":"cid-1"}`
	hookStdinJSON, err := ReadHookStdinJSON(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadHookStdinJSON() error = %v", err)
	}
	if string(hookStdinJSON) != input {
		t.Fatalf("unexpected body: got %q, want %q", hookStdinJSON, input)
	}
}

func TestReadHookStdinJSON_StripsUTF8BOM(t *testing.T) {
	input := "\uFEFF" + `{"hook_event_name":"afterFileEdit","conversation_id":"cid-1"}`
	hookStdinJSON, err := ReadHookStdinJSON(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadHookStdinJSON() error = %v", err)
	}
	if string(hookStdinJSON) != `{"hook_event_name":"afterFileEdit","conversation_id":"cid-1"}` {
		t.Fatalf("unexpected body: %q", hookStdinJSON)
	}
}

func TestReadHookStdinJSON_InvalidJSON(t *testing.T) {
	_, err := ReadHookStdinJSON(strings.NewReader(`not json`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Fatalf("ReadHookStdinJSON() error = %v, want message containing %q", err, "invalid JSON")
	}
}
