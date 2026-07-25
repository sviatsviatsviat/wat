package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	watBinaryOnce sync.Once
	watBinaryPath string
	watBinaryErr  error
)

func buildWat(t *testing.T) string {
	t.Helper()

	watBinaryOnce.Do(func() {
		root := moduleRoot()
		dir, err := os.MkdirTemp("", "wat-e2e-bin-*")
		if err != nil {
			watBinaryErr = err
			return
		}
		name := "wat"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		watBinaryPath = filepath.Join(dir, name)
		cmd := exec.Command("go", "build", "-buildvcs=true", "-o", watBinaryPath, "./cmd/wat")
		cmd.Dir = root
		var out []byte
		out, watBinaryErr = cmd.CombinedOutput()
		if watBinaryErr != nil {
			watBinaryErr = &buildError{err: watBinaryErr, out: out}
		}
	})
	if watBinaryErr != nil {
		t.Fatalf("go build ./cmd/wat: %v", watBinaryErr)
	}
	return watBinaryPath
}

type buildError struct {
	err error
	out []byte
}

func (e *buildError) Error() string {
	return e.err.Error() + "\n" + string(e.out)
}

func moduleRoot() string {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		panic("go env GOMOD: " + err.Error())
	}
	mod := strings.TrimSpace(string(out))
	if mod == "" || mod == "/dev/null" {
		panic("not inside a Go module")
	}
	return filepath.Dir(mod)
}

func watModuleVersion(t *testing.T, binary string) string {
	t.Helper()

	cmd := exec.Command("go", "version", "-m", binary)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go version -m: %v\n%s", err, out)
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "mod\t") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[2] != "" && fields[2] != "(devel)" {
			return fields[2]
		}
	}
	// Local builds often report (devel); init still needs a require version,
	// and the replace directive below points at the checkout.
	return "v0.0.0-e2e-000000000000"
}

// initProjectWithReplace scaffolds .wat/ with wat init after seeding go.mod with
// a replace directive so tidy resolves the local checkout.
func initProjectWithReplace(t *testing.T) string {
	t.Helper()

	binary := buildWat(t)
	version := watModuleVersion(t, binary)
	dir := t.TempDir()
	watDir := filepath.Join(dir, ".wat")
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}

	goMod := "module wat-hooks\n\ngo 1.26\n\nrequire github.com/sviatsviatsviat/wat " + version + "\n\nreplace github.com/sviatsviatsviat/wat => " + filepath.ToSlash(moduleRoot()) + "\n"
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, stderr, code := runWat(t, binary, dir, "init")
	if code != 0 {
		t.Fatalf("wat init exit = %d\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}
	return dir
}

func runWat(t *testing.T, binary, dir string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"WAT_PROJECT_DIR="+dir,
		"PATH="+filepath.Dir(binary)+string(os.PathListSeparator)+os.Getenv("PATH"),
	)

	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("run wat %v: %v\nstdout:\n%s\nstderr:\n%s", args, err, outBuf.String(), errBuf.String())
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func fixturePath(t *testing.T, agent, name string) string {
	t.Helper()
	return filepath.Join(moduleRoot(), "testdata", "fixtures", agent, name)
}
