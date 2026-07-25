package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/initproj"
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
	deps := initproj.DefaultDeps()
	deps.Command = func(name string, args ...string) *exec.Cmd {
		gotCmdName = name
		gotCmdArgs = append([]string(nil), args...)
		gotCmd = exec.Command("go", "version")
		return gotCmd
	}

	if err := initproj.Init(dir, false, watModuleVersionFn(), deps, stdout, stderr); err != nil {
		t.Fatalf("initproj.Init: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".wat")); err != nil {
		t.Fatalf("missing .wat/: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".wat", ".cache")); err != nil {
		t.Fatalf("missing .wat/.cache/: %v", err)
	}
	hooksPath := filepath.Join(dir, ".wat", "hooks.go")
	hooks, err := os.ReadFile(hooksPath)
	if err != nil {
		t.Fatalf("read hooks.go: %v", err)
	}
	if !strings.Contains(string(hooks), "agnostic.UseHooks()") {
		t.Fatalf("hooks.go missing UseHooks: %s", hooks)
	}
	goModPath := filepath.Join(dir, ".wat", "go.mod")
	goMod, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "github.com/sviatsviatsviat/wat v0.0.0-test-000000000000") {
		t.Fatalf("go.mod missing require: %s", goMod)
	}
	if gotCmdName != "go" || len(gotCmdArgs) != 2 || gotCmdArgs[0] != "mod" || gotCmdArgs[1] != "tidy" {
		t.Fatalf("expected go mod tidy, got %s %v", gotCmdName, gotCmdArgs)
	}
	if gotCmd == nil || gotCmd.Dir != filepath.Join(dir, ".wat") {
		t.Fatalf("go mod tidy Dir = %v", gotCmd)
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
	deps := initproj.DefaultDeps()
	deps.Command = fakeCommand

	if err := initproj.Init(dir, false, watModuleVersionFn(), deps, stdout, stderr); err != nil {
		t.Fatalf("initproj.Init first run: %v", err)
	}

	err := initproj.Init(dir, false, watModuleVersionFn(), deps, stdout, stderr)
	if err == nil {
		t.Fatalf("initproj.Init second run: expected error")
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Fatalf("error = %q, want --force mention", err)
	}

	if err := initproj.Init(dir, true, watModuleVersionFn(), deps, stdout, stderr); err != nil {
		t.Fatalf("initproj.Init with force: %v", err)
	}
}
