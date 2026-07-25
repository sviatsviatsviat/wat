package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunTest_missingFixtureUsage(t *testing.T) {
	prevStderr := stderr
	var errBuf bytes.Buffer
	stderr = &errBuf
	t.Cleanup(func() { stderr = prevStderr })

	code := run([]string{"test"})
	if code != exitUsage {
		t.Fatalf("exit = %d, want %d", code, exitUsage)
	}
	if !strings.Contains(errBuf.String(), "--fixture is required") {
		t.Fatalf("stderr = %q", errBuf.String())
	}
}
