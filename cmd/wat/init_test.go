package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitProject_createsWatDirAndFiles(t *testing.T) {
	prevVersionFn := watModuleVersionFn
	watModuleVersionFn = func() string { return "v0.0.0-test-000000000000" }
	t.Cleanup(func() { watModuleVersionFn = prevVersionFn })

	dir := t.TempDir()
	prevStdout, prevStderr := stdout, stderr
	stdout, stderr = io.Discard, io.Discard
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	var gotCmdName string
	var gotCmdArgs []string
	var gotCmd *exec.Cmd
	fakeCommand := func(name string, args ...string) *exec.Cmd {
		gotCmdName = name
		gotCmdArgs = append([]string(nil), args...)
		gotCmd = exec.Command("go", "version")
		return gotCmd
	}

	if err := initProject(dir, false, fakeCommand); err != nil {
		t.Fatalf("initProject: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".wat")); err != nil {
		t.Fatalf("missing .wat/: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".wat", ".cache")); err != nil {
		t.Fatalf("missing .wat/.cache/: %v", err)
	}

	goModBytes, err := os.ReadFile(filepath.Join(dir, ".wat", "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if !strings.Contains(string(goModBytes), "go 1.26") {
		t.Fatalf("go.mod missing go directive: %q", string(goModBytes))
	}

	hooksBytes, err := os.ReadFile(filepath.Join(dir, ".wat", "hooks.go"))
	if err != nil {
		t.Fatalf("read hooks.go: %v", err)
	}
	if !strings.Contains(string(hooksBytes), "agenthooks.NewMux") {
		t.Fatalf("hooks.go missing mux usage: %q", string(hooksBytes))
	}
	if !strings.Contains(string(hooksBytes), "github.com/sviatsviatsviat/wat/agenthooks") {
		t.Fatalf("hooks.go missing agenthooks import: %q", string(hooksBytes))
	}

	if gotCmdName != "go" || len(gotCmdArgs) != 2 || gotCmdArgs[0] != "mod" || gotCmdArgs[1] != "tidy" {
		t.Fatalf("go mod tidy invocation = %q %q, want %q", gotCmdName, gotCmdArgs, []string{"mod", "tidy"})
	}
	if gotCmd == nil {
		t.Fatalf("expected go mod tidy command to be created")
	}
	if gotCmd.Dir != filepath.Join(dir, ".wat") {
		t.Fatalf("go mod tidy dir = %q, want %q", gotCmd.Dir, filepath.Join(dir, ".wat"))
	}
}

func TestInitProject_refusesOverwriteWithoutForce(t *testing.T) {
	prevVersionFn := watModuleVersionFn
	watModuleVersionFn = func() string { return "v0.0.0-test-000000000000" }
	t.Cleanup(func() { watModuleVersionFn = prevVersionFn })

	dir := t.TempDir()
	prevStdout, prevStderr := stdout, stderr
	stdout, stderr = io.Discard, io.Discard
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	fakeCommand := func(string, ...string) *exec.Cmd {
		return exec.Command("go", "version")
	}

	if err := initProject(dir, false, fakeCommand); err != nil {
		t.Fatalf("initProject first run: %v", err)
	}

	sentinel := []byte("package main\n\nconst sentinel = true\n")
	hooksPath := filepath.Join(dir, ".wat", "hooks.go")
	if err := os.WriteFile(hooksPath, sentinel, 0o644); err != nil {
		t.Fatalf("write sentinel hooks.go: %v", err)
	}

	if err := initProject(dir, false, fakeCommand); err == nil {
		t.Fatalf("initProject second run: expected error")
	}
	got, err := os.ReadFile(hooksPath)
	if err != nil {
		t.Fatalf("read hooks.go: %v", err)
	}
	if string(got) != string(sentinel) {
		t.Fatalf("hooks.go changed without --force: got %q want %q", string(got), string(sentinel))
	}

	if err := initProject(dir, true, fakeCommand); err != nil {
		t.Fatalf("initProject with force: %v", err)
	}
	got, err = os.ReadFile(hooksPath)
	if err != nil {
		t.Fatalf("read hooks.go after --force: %v", err)
	}
	if string(got) != watHooksGo {
		t.Fatalf("hooks.go content after --force mismatch")
	}
}

