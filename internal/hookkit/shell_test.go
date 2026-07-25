package hookkit

import (
	"encoding/json"
	"testing"
)

func TestExtractShellCommand(t *testing.T) {
	t.Parallel()
	if got := ExtractShellCommand(json.RawMessage(`{"command":"echo hi"}`)); got != "echo hi" {
		t.Fatalf("ExtractShellCommand() = %q", got)
	}
	if got := ExtractShellCommand(nil); got != "" {
		t.Fatalf("ExtractShellCommand(nil) = %q", got)
	}
}
