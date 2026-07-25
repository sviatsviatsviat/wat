package dotwat_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// buildWat compiles the wat CLI into a temp directory and returns its path.
// -buildvcs=true is required: watModuleVersionFn reads VCS build info to derive
// the version used for the .wat/ content-addressed build cache key.
func buildWat(t *testing.T) string {
	t.Helper()

	root := moduleRoot(t)
	dir := t.TempDir()
	name := "wat"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binPath := filepath.Join(dir, name)

	cmd := exec.Command("go", "build", "-buildvcs=true", "-o", binPath, "./cmd/wat")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./cmd/wat: %v\n%s", err, out)
	}
	return binPath
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		t.Fatalf("go env GOMOD: %v", err)
	}
	mod := strings.TrimSpace(string(out))
	if mod == "" || mod == "/dev/null" {
		t.Fatal("not inside a Go module")
	}
	return filepath.Dir(mod)
}

// runWatTest runs "wat test --agent cursor --fixture <fixture>" against this
// repository's own .wat/ project and returns combined stdout+stderr and exit code.
func runWatTest(t *testing.T, binary, root, fixture string) (output string, exitCode int) {
	t.Helper()

	cmd := exec.Command(binary, "test", "--agent", "cursor", "--fixture", fixture)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "WAT_PROJECT_DIR="+root)

	var buf strings.Builder
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	if err == nil {
		return buf.String(), 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("run wat test: %v\noutput:\n%s", err, buf.String())
	}
	return buf.String(), exitErr.ExitCode()
}

// TestCursorSubagentStart_ModelGate runs this repository's committed .wat/hooks.go
// against real cursor subagentStart fixtures, verifying the gate that denies a
// subagent spawn pinned to a non-"auto" model and allows the "auto" model through.
func TestCursorSubagentStart_ModelGate(t *testing.T) {
	root := moduleRoot(t)
	binary := buildWat(t)

	t.Run("non-auto model is denied", func(t *testing.T) {
		fixture := filepath.Join(root, "testdata", "fixtures", "cursor", "subagent_start.json")
		out, code := runWatTest(t, binary, root, fixture)
		if code != 2 {
			t.Fatalf("exit = %d, want 2\noutput:\n%s", code, out)
		}
		if !strings.Contains(out, `"permission":"deny"`) {
			t.Fatalf("output missing deny permission:\n%s", out)
		}
		if !strings.Contains(out, "Re-run it with the auto model") {
			t.Fatalf("output missing user message:\n%s", out)
		}
	})

	t.Run("auto model is allowed", func(t *testing.T) {
		fixture := filepath.Join(root, "testdata", "fixtures", "cursor", "subagent_start_auto.json")
		out, code := runWatTest(t, binary, root, fixture)
		if code != 0 {
			t.Fatalf("exit = %d, want 0\noutput:\n%s", code, out)
		}
		if !strings.Contains(out, "stdout: (empty)") {
			t.Fatalf("output missing empty stdout marker:\n%s", out)
		}
	})
}
