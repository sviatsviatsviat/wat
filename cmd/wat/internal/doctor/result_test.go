package doctor

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintResult_knownStatuses(t *testing.T) {
	var buf bytes.Buffer
	PrintResult(&buf, Result{Group: "cache", Status: Pass, Message: "ok"})
	if got := buf.String(); !strings.HasPrefix(got, "PASS") {
		t.Fatalf("pass prefix = %q", got)
	}
	buf.Reset()
	PrintResult(&buf, Result{Group: "cache", Status: Fail, Message: "bad", Fix: "fix it"})
	if got := buf.String(); !strings.Contains(got, "FAIL") || !strings.Contains(got, "fix: fix it") {
		t.Fatalf("fail output = %q", got)
	}
	buf.Reset()
	PrintResult(&buf, Result{Group: "cache", Status: Warn, Message: "maybe"})
	if got := buf.String(); !strings.HasPrefix(got, "WARN") {
		t.Fatalf("warn prefix = %q", got)
	}
}

func TestPrintResult_unknownStatus(t *testing.T) {
	var buf bytes.Buffer
	PrintResult(&buf, Result{Group: "cache", Status: Status(99), Message: "odd"})
	if got := buf.String(); !strings.HasPrefix(got, "????") {
		t.Fatalf("unknown prefix = %q", got)
	}
}
