// Package proctest provides cross-platform subprocess stubs for CLI tests.
//
// Import this package only from *_test.go files. When a child process is
// started with WAT_PROCTEST_MODE set, init exits before tests run so the
// current test binary can stand in for Unix utilities such as echo and sleep.
package proctest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	envMode = "WAT_PROCTEST_MODE"
	envOut  = "WAT_PROCTEST_OUT"
	envPath = "WAT_PROCTEST_PATH"
)

func init() {
	switch os.Getenv(envMode) {
	case "println":
		fmt.Println(os.Getenv(envOut))
		os.Exit(0)
	case "print":
		fmt.Print(os.Getenv(envOut))
		os.Exit(0)
	case "stderr_exit1":
		fmt.Fprintln(os.Stderr, os.Getenv(envOut))
		os.Exit(1)
	case "exit0":
		os.Exit(0)
	case "sleep":
		time.Sleep(time.Hour)
		os.Exit(0)
	case "touch":
		path := os.Getenv(envPath)
		_ = os.MkdirAll(filepath.Dir(path), 0o750) //nolint:gosec // test-only path from WAT_PROCTEST_PATH
		_ = os.WriteFile(path, nil, 0o600)         //nolint:gosec // test-only path from WAT_PROCTEST_PATH
		os.Exit(0)
	}
}

// Println returns a command that prints s followed by a newline and exits 0.
func Println(s string) *exec.Cmd {
	return cmd("println", envOut+"="+s)
}

// ExitZero returns a command that exits 0 with no output.
func ExitZero() *exec.Cmd {
	return cmd("exit0")
}

// Fail returns a command that prints msg to stderr and exits 1.
func Fail(msg string) *exec.Cmd {
	return cmd("stderr_exit1", envOut+"="+msg)
}

// Sleep returns a command that blocks until killed (for timeout tests).
func Sleep() *exec.Cmd {
	return cmd("sleep")
}

// Touch returns a command that creates path (and parents) then exits 0.
func Touch(path string) *exec.Cmd {
	return cmd("touch", envPath+"="+path)
}

func cmd(mode string, extraEnv ...string) *exec.Cmd {
	bin := os.Args[0]
	if abs, err := filepath.Abs(bin); err == nil {
		bin = abs
	}
	c := exec.Command(bin, "-test.run=^$", "-test.count=1") //nolint:gosec // re-exec current test binary as a stub
	env := append(os.Environ(), envMode+"="+mode)
	env = append(env, extraEnv...)
	c.Env = env
	return c
}
