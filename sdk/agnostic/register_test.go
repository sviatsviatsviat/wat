package agnostic

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestOnNilIgnored(t *testing.T) {
	resetTest(t)
	OnPreTool(nil)
	OnAny(nil)

	var stdout bytes.Buffer
	code := Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatal("nil handlers should not register")
	}
}
