package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersion_printsModuleVersion(t *testing.T) {
	prevStdout, prevStderr := stdout, stderr
	prevVersionFn := watModuleVersionFn
	var outBuf, errBuf bytes.Buffer
	stdout, stderr = &outBuf, &errBuf
	watModuleVersionFn = func() string { return "v0.1.1-alpha" }
	t.Cleanup(func() {
		stdout, stderr = prevStdout, prevStderr
		watModuleVersionFn = prevVersionFn
	})

	code := run([]string{"version"})
	if code != exitOK {
		t.Fatalf("exit = %d, want %d", code, exitOK)
	}
	if got := strings.TrimSpace(outBuf.String()); got != "v0.1.1-alpha" {
		t.Fatalf("stdout = %q, want %q", got, "v0.1.1-alpha")
	}
	if errBuf.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", errBuf.String())
	}
}

func TestRunVersion_rootFlags(t *testing.T) {
	prevStdout, prevStderr := stdout, stderr
	prevVersionFn := watModuleVersionFn
	watModuleVersionFn = func() string { return "v0.0.0-test-000000000000" }
	t.Cleanup(func() {
		stdout, stderr = prevStdout, prevStderr
		watModuleVersionFn = prevVersionFn
	})

	for _, args := range [][]string{{"-version"}, {"--version"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer
			stdout, stderr = &outBuf, &errBuf

			code := run(args)
			if code != exitOK {
				t.Fatalf("exit = %d, want %d", code, exitOK)
			}
			if got := strings.TrimSpace(outBuf.String()); got != "v0.0.0-test-000000000000" {
				t.Fatalf("stdout = %q, want %q", got, "v0.0.0-test-000000000000")
			}
			if errBuf.Len() != 0 {
				t.Fatalf("stderr = %q, want empty", errBuf.String())
			}
		})
	}
}

func TestRunVersion_missingBuildInfo(t *testing.T) {
	prevStdout, prevStderr := stdout, stderr
	prevVersionFn := watModuleVersionFn
	var outBuf, errBuf bytes.Buffer
	stdout, stderr = &outBuf, &errBuf
	watModuleVersionFn = func() string { return "" }
	t.Cleanup(func() {
		stdout, stderr = prevStdout, prevStderr
		watModuleVersionFn = prevVersionFn
	})

	code := run([]string{"version"})
	if code != exitRuntimeFailure {
		t.Fatalf("exit = %d, want %d", code, exitRuntimeFailure)
	}
	if outBuf.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", outBuf.String())
	}
	if !strings.Contains(errBuf.String(), "determine wat module version") {
		t.Fatalf("stderr = %q", errBuf.String())
	}
}

func TestRunVersion_help(t *testing.T) {
	prevStdout, prevStderr := stdout, stderr
	var outBuf, errBuf bytes.Buffer
	stdout, stderr = &outBuf, &errBuf
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	code := run([]string{"version", "-h"})
	if code != exitOK {
		t.Fatalf("exit = %d, want %d", code, exitOK)
	}
	if !strings.Contains(outBuf.String(), "wat version") {
		t.Fatalf("stdout missing help:\n%s", outBuf.String())
	}
	if errBuf.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", errBuf.String())
	}
}
