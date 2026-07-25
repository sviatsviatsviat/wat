package e2e_test

import (
	"strings"
	"testing"
)

func TestWatDoctor_allPass(t *testing.T) {
	binary := buildWat(t)
	project := initProjectWithReplace(t)

	stdout, stderr, code := runWat(t, binary, project, "install", "--wat-path", binary)
	if code != 0 {
		t.Fatalf("install exit = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}

	stdout, stderr, code = runWat(t, binary, project, "doctor")
	if code != 0 {
		t.Fatalf("doctor exit = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}
	out := stdout + stderr
	for _, want := range []string{"PASS", "toolchain", "script", "cache", "install"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "FAIL") {
		t.Fatalf("unexpected FAIL in output:\n%s", out)
	}
}
