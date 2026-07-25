package buildcache

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type fakeFileInfo struct{}

func (fakeFileInfo) Name() string       { return "x" }
func (fakeFileInfo) Size() int64        { return 1 }
func (fakeFileInfo) Mode() os.FileMode  { return 0o644 }
func (fakeFileInfo) ModTime() time.Time { return time.Now() }
func (fakeFileInfo) IsDir() bool        { return false }
func (fakeFileInfo) Sys() any           { return nil }

type fakeDirEntry struct {
	name  string
	isDir bool
}

func (e fakeDirEntry) Name() string               { return e.name }
func (e fakeDirEntry) IsDir() bool                { return e.isDir }
func (e fakeDirEntry) Type() os.FileMode          { return 0 }
func (e fakeDirEntry) Info() (os.FileInfo, error) { return fakeFileInfo{}, nil }

func TestCacheKey_changesOnGoMod(t *testing.T) {
	deps := DefaultDeps()
	deps.Getenv = func(string) string { return "" }
	files := map[string][]byte{
		"go.mod": []byte("mod"),
		"go.sum": []byte("sum"),
		"a.go":   []byte("hooks-a"),
		"b.go":   []byte("hooks-b"),
	}
	deps.ReadDir = func(dir string) ([]os.DirEntry, error) {
		if filepath.Base(dir) != ".wat" && dir != "/tmp/.wat" {
			return nil, nil
		}
		return []os.DirEntry{
			fakeDirEntry{name: "go.mod"},
			fakeDirEntry{name: "go.sum"},
			fakeDirEntry{name: "a.go"},
			fakeDirEntry{name: "b.go"},
		}, nil
	}
	deps.ReadFile = func(path string) ([]byte, error) {
		if data, ok := files[filepath.Base(path)]; ok {
			return data, nil
		}
		return nil, os.ErrNotExist
	}

	a, err := CacheKey("/tmp/.wat", "v1", deps)
	if err != nil {
		t.Fatal(err)
	}
	files["go.mod"] = []byte("mod2")
	b, err := CacheKey("/tmp/.wat", "v1", deps)
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatalf("cache key should differ when go.mod changes")
	}
}

func TestCacheKey_includesNestedFiles(t *testing.T) {
	deps := DefaultDeps()
	deps.Getenv = func(string) string { return "" }
	deps.ReadDir = func(dir string) ([]os.DirEntry, error) {
		switch filepath.ToSlash(dir) {
		case "/tmp/.wat":
			return []os.DirEntry{
				fakeDirEntry{name: "go.mod"},
				fakeDirEntry{name: "hooks.go"},
				fakeDirEntry{name: "internal", isDir: true},
				fakeDirEntry{name: ".cache", isDir: true},
			}, nil
		case "/tmp/.wat/internal":
			return []os.DirEntry{fakeDirEntry{name: "lib.go"}}, nil
		case "/tmp/.wat/.cache":
			t.Fatal(".cache should be skipped")
			return nil, nil
		default:
			return nil, os.ErrNotExist
		}
	}
	deps.ReadFile = func(path string) ([]byte, error) {
		switch filepath.ToSlash(path) {
		case "/tmp/.wat/go.mod":
			return []byte("module x\n"), nil
		case "/tmp/.wat/hooks.go":
			return []byte("package main\n"), nil
		case "/tmp/.wat/internal/lib.go":
			return []byte("package internal\n"), nil
		default:
			return nil, os.ErrNotExist
		}
	}

	withNested, err := CacheKey("/tmp/.wat", "v1", deps)
	if err != nil {
		t.Fatal(err)
	}

	deps.ReadFile = func(path string) ([]byte, error) {
		switch filepath.ToSlash(path) {
		case "/tmp/.wat/go.mod":
			return []byte("module x\n"), nil
		case "/tmp/.wat/hooks.go":
			return []byte("package main\n"), nil
		case "/tmp/.wat/internal/lib.go":
			return []byte("package internal\n// changed\n"), nil
		default:
			return nil, os.ErrNotExist
		}
	}
	changed, err := CacheKey("/tmp/.wat", "v1", deps)
	if err != nil {
		t.Fatal(err)
	}
	if withNested == changed {
		t.Fatal("cache key should change when nested package file changes")
	}
}

func TestCacheKey_includesBuildEnv(t *testing.T) {
	deps := DefaultDeps()
	deps.Getenv = func(key string) string {
		if key == "CGO_ENABLED" {
			return "0"
		}
		return ""
	}
	deps.ReadDir = func(string) ([]os.DirEntry, error) {
		return []os.DirEntry{fakeDirEntry{name: "hooks.go"}}, nil
	}
	deps.ReadFile = func(string) ([]byte, error) { return []byte("package main\n"), nil }

	a, err := CacheKey("/tmp/.wat", "v1", deps)
	if err != nil {
		t.Fatal(err)
	}
	deps.Getenv = func(key string) string {
		if key == "CGO_ENABLED" {
			return "1"
		}
		return ""
	}
	b, err := CacheKey("/tmp/.wat", "v1", deps)
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("cache key should differ when CGO_ENABLED changes")
	}
}

func TestCurrentModulePath_rejectsMultiLine(t *testing.T) {
	dir := t.TempDir()
	deps := DefaultDeps()
	deps.Command = func(name string, args ...string) *exec.Cmd {
		if name == "go" && len(args) >= 1 && args[0] == "list" {
			return exec.Command("printf", "%s", "github.com/example/a\nwat-hooks\n")
		}
		return exec.Command(name, args...)
	}
	_, err := currentModulePath(dir, deps, buildSettings{})
	if err == nil || !strings.Contains(err.Error(), "multi-line") {
		t.Fatalf("got err %v, want multi-line failure", err)
	}
}

func TestPinnedBuildEnv_overrides(t *testing.T) {
	got := pinnedBuildEnv([]string{"GOOS=windows", "PATH=/bin"}, buildSettings{
		goos:       "linux",
		goarch:     "amd64",
		goflags:    "-trimpath",
		cgoEnabled: "0",
	})
	want := map[string]string{
		"GOOS":        "linux",
		"GOARCH":      "amd64",
		"GOFLAGS":     "-trimpath",
		"CGO_ENABLED": "0",
		"PATH":        "/bin",
	}
	for _, e := range got {
		k, v, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		if want[k] != v {
			t.Fatalf("%s = %q, want %q", k, v, want[k])
		}
		delete(want, k)
	}
	if len(want) != 0 {
		t.Fatalf("missing env keys: %v", want)
	}
}
