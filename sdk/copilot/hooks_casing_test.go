package copilot

import (
	"strings"
	"testing"
)

func TestPackageCodecRejectsCamelCaseEvents(t *testing.T) {
	t.Parallel()
	_, err := codec.DecodeAs([]byte(`{}`), "preToolUse")
	if err == nil {
		t.Fatal("expected decode error for camelCase event name")
	}
	if !strings.Contains(err.Error(), `use PascalCase "PreToolUse"`) {
		t.Fatalf("err = %v", err)
	}
}
