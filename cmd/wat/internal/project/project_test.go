package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindFrom_walkUp(t *testing.T) {
	root := t.TempDir()
	watDir := Dir(root)
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, HooksFile), []byte("package main\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, GoModFile), []byte("module x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	deps := DefaultDeps()
	got, err := FindFrom(nested, deps)
	if err != nil {
		t.Fatal(err)
	}
	if got != watDir {
		t.Fatalf("FindFrom = %q, want %q", got, watDir)
	}
}

func TestResolve_envOverrideErrors(t *testing.T) {
	deps := DefaultDeps()
	deps.Getenv = func(key string) string {
		if key == ProjectDirEnv {
			return "/does/not/exist"
		}
		return ""
	}
	deps.Stat = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }

	_, err := Resolve(deps)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), DirName) && !strings.Contains(err.Error(), ProjectDirEnv) {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestMustHaveFiles_missing(t *testing.T) {
	dir := t.TempDir()
	err := MustHaveFiles(dir, DefaultDeps())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMustHaveFiles_rejectsDirectory(t *testing.T) {
	dir := t.TempDir()
	watDir := Dir(dir)
	if err := os.MkdirAll(filepath.Join(watDir, HooksFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, GoModFile), []byte("module x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := MustHaveFiles(watDir, DefaultDeps())
	if err == nil {
		t.Fatal("expected error for non-regular hooks.go")
	}
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("error = %q", err)
	}
}

func TestDir(t *testing.T) {
	if got := Dir("/project"); got != filepath.Join("/project", DirName) {
		t.Fatalf("Dir = %q", got)
	}
}
