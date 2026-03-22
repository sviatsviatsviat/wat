package cli

import (
	"strings"
	"testing"
)

func TestMockConsole_recordsWritesAndStderrBufferWriter(t *testing.T) {
	t.Parallel()
	mockConsole := NewMockConsole()
	_ = mockConsole.WriteError("err line")
	_, _ = mockConsole.StderrBufferWriter().Write([]byte("child\n"))

	recordedStderr := mockConsole.StderrString()
	if !strings.Contains(recordedStderr, "err line") || !strings.Contains(recordedStderr, "child") {
		t.Fatalf("StderrString() = %q, want both err line and child", recordedStderr)
	}
	if !mockConsole.StderrContains("child") {
		t.Fatal("StderrContains should find child output")
	}
}

func TestMockConsole_Write_goesToStdoutBuffer(t *testing.T) {
	t.Parallel()
	mockConsole := NewMockConsole()
	if err := mockConsole.Write("{}\n"); err != nil {
		t.Fatalf("Write: %v", err)
	}
	_, _ = mockConsole.StdoutBufferWriter().Write([]byte("extra"))
	if want := "{}\nextra"; mockConsole.StdoutString() != want {
		t.Fatalf("StdoutString: got %q want %q", mockConsole.StdoutString(), want)
	}
	if !mockConsole.StdoutContains("{}") {
		t.Fatal("StdoutContains should find default response body")
	}
	if mockConsole.StderrString() != "" {
		t.Fatalf("stderr buffer should be empty, got %q", mockConsole.StderrString())
	}
}
