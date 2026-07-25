package doctor

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintResult_knownStatuses(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
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
	t.Setenv("NO_COLOR", "1")
	var buf bytes.Buffer
	PrintResult(&buf, Result{Group: "cache", Status: Status(99), Message: "odd"})
	if got := buf.String(); !strings.HasPrefix(got, "????") {
		t.Fatalf("unknown prefix = %q", got)
	}
}

func TestPrintResult_forceColor(t *testing.T) {
	t.Setenv("FORCE_COLOR", "1")
	t.Setenv("NO_COLOR", "")
	var buf bytes.Buffer
	PrintResult(&buf, Result{Group: "cache", Status: Pass, Message: "ok"})
	got := buf.String()
	if !strings.Contains(got, ansiGreen) || !strings.Contains(got, ansiReset) {
		t.Fatalf("expected green PASS, got %q", got)
	}
	buf.Reset()
	PrintResult(&buf, Result{Group: "install", Status: Warn, Message: "gap"})
	got = buf.String()
	if !strings.Contains(got, ansiYellow) {
		t.Fatalf("expected yellow WARN, got %q", got)
	}
	buf.Reset()
	PrintResult(&buf, Result{Group: "script", Status: Fail, Message: "bad"})
	got = buf.String()
	if !strings.Contains(got, ansiRed) {
		t.Fatalf("expected red FAIL, got %q", got)
	}
}

func TestPrintResult_noColorOverridesForce(t *testing.T) {
	t.Setenv("FORCE_COLOR", "1")
	t.Setenv("NO_COLOR", "1")
	var buf bytes.Buffer
	PrintResult(&buf, Result{Group: "cache", Status: Pass, Message: "ok"})
	if got := buf.String(); strings.Contains(got, "\033[") {
		t.Fatalf("NO_COLOR should disable color, got %q", got)
	}
}

func TestPrintFailureSummary(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var buf bytes.Buffer
	PrintFailureSummary(&buf, 0)
	if buf.Len() != 0 {
		t.Fatalf("expected empty summary for 0 fails, got %q", buf.String())
	}
	PrintFailureSummary(&buf, 2)
	if got := buf.String(); got != "\nwat doctor: 2 check(s) failed\n" {
		t.Fatalf("summary = %q", got)
	}
}
