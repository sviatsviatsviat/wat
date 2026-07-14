package adapter

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

func TestNewToolCall_BashShell(t *testing.T) {
	tc := NewToolCall("Bash", []byte(`{"command":"echo hi","description":"test"}`), "id1")
	if tc.Name != model.ToolBash || tc.Native != "Bash" || tc.Shell != "echo hi" || tc.ID != "id1" {
		t.Fatalf("ToolCall=%+v", tc)
	}
}

func TestRawToText(t *testing.T) {
	if got := RawToText([]byte(`"hello"`)); got != "hello" {
		t.Fatalf("string = %q", got)
	}
	if got := RawToText([]byte(`{"k":"v"}`)); got != `{"k":"v"}` {
		t.Fatalf("object = %q", got)
	}
}
