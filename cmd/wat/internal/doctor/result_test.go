package doctor

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestPrintResult_knownStatuses(t *testing.T) {
	var buf bytes.Buffer
	PrintResult(&buf, Result{Group: "cache", Status: Pass, Message: "ok"}, false)
	if got := buf.String(); !strings.HasPrefix(got, "PASS") {
		t.Fatalf("pass prefix = %q", got)
	}
	buf.Reset()
	PrintResult(&buf, Result{Group: "cache", Status: Fail, Message: "bad", Fix: "fix it"}, false)
	if got := buf.String(); !strings.Contains(got, "FAIL") || !strings.Contains(got, "fix: fix it") {
		t.Fatalf("fail output = %q", got)
	}
	buf.Reset()
	PrintResult(&buf, Result{Group: "cache", Status: Warn, Message: "maybe"}, false)
	if got := buf.String(); !strings.HasPrefix(got, "WARN") {
		t.Fatalf("warn prefix = %q", got)
	}
}

func TestPrintResult_unknownStatus(t *testing.T) {
	var buf bytes.Buffer
	PrintResult(&buf, Result{Group: "cache", Status: Status(99), Message: "odd"}, false)
	if got := buf.String(); !strings.HasPrefix(got, "????") {
		t.Fatalf("unknown prefix = %q", got)
	}
}

func TestPrintResult_color(t *testing.T) {
	var buf bytes.Buffer
	PrintResult(&buf, Result{Group: "cache", Status: Pass, Message: "ok"}, true)
	got := buf.String()
	if !strings.Contains(got, ansiGreen) || !strings.Contains(got, ansiReset) {
		t.Fatalf("expected green PASS, got %q", got)
	}
	buf.Reset()
	PrintResult(&buf, Result{Group: "install", Status: Warn, Message: "gap"}, true)
	got = buf.String()
	if !strings.Contains(got, ansiYellow) {
		t.Fatalf("expected yellow WARN, got %q", got)
	}
	buf.Reset()
	PrintResult(&buf, Result{Group: "script", Status: Fail, Message: "bad"}, true)
	got = buf.String()
	if !strings.Contains(got, ansiRed) {
		t.Fatalf("expected red FAIL, got %q", got)
	}
}

func TestColorEnabled(t *testing.T) {
	t.Run("no_color", func(t *testing.T) {
		deps := Deps{
			LookupEnv: func(key string) (string, bool) {
				if key == "NO_COLOR" {
					return "1", true
				}
				return "", false
			},
			Getenv:     func(string) string { return "1" },
			IsTerminal: func(io.Writer) bool { return true },
		}
		if ColorEnabled(deps, io.Discard) {
			t.Fatal("NO_COLOR should disable color")
		}
	})
	t.Run("force_color", func(t *testing.T) {
		deps := Deps{
			LookupEnv: func(string) (string, bool) { return "", false },
			Getenv: func(key string) string {
				if key == "FORCE_COLOR" {
					return "1"
				}
				return ""
			},
			IsTerminal: func(io.Writer) bool { return false },
		}
		if !ColorEnabled(deps, io.Discard) {
			t.Fatal("FORCE_COLOR should enable color")
		}
	})
	t.Run("terminal", func(t *testing.T) {
		deps := Deps{
			LookupEnv:  func(string) (string, bool) { return "", false },
			Getenv:     func(string) string { return "" },
			IsTerminal: func(io.Writer) bool { return true },
		}
		if !ColorEnabled(deps, io.Discard) {
			t.Fatal("terminal should enable color")
		}
	})
}

func TestPrintFailureSummary(t *testing.T) {
	var buf bytes.Buffer
	PrintFailureSummary(&buf, 0, false)
	if buf.Len() != 0 {
		t.Fatalf("expected empty summary for 0 fails, got %q", buf.String())
	}
	PrintFailureSummary(&buf, 2, false)
	if got := buf.String(); got != "\nwat doctor: 2 check(s) failed\n" {
		t.Fatalf("summary = %q", got)
	}

	buf.Reset()
	PrintFailureSummary(&buf, 2, true)
	got := buf.String()
	want := "\n" + ansiRed + "wat doctor: 2 check(s) failed" + ansiReset + "\n"
	if got != want {
		t.Fatalf("colored summary = %q, want %q", got, want)
	}
}
