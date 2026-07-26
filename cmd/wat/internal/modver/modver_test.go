package modver

import (
	"runtime/debug"
	"testing"
	"time"
)

func TestResolve_override(t *testing.T) {
	prevOverride := Override
	prevRead := readBuildInfo
	Override = "  v0.1.1-alpha  "
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		t.Fatal("readBuildInfo should not run when Override is set")
		return nil, false
	}
	t.Cleanup(func() {
		Override = prevOverride
		readBuildInfo = prevRead
	})

	if got := Resolve(); got != "v0.1.1-alpha" {
		t.Fatalf("Resolve() = %q, want %q", got, "v0.1.1-alpha")
	}
}

func TestResolve_moduleVersion(t *testing.T) {
	prevOverride := Override
	prevRead := readBuildInfo
	Override = ""
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Main: debug.Module{Version: "v0.2.0"}}, true
	}
	t.Cleanup(func() {
		Override = prevOverride
		readBuildInfo = prevRead
	})

	if got := Resolve(); got != "v0.2.0" {
		t.Fatalf("Resolve() = %q, want %q", got, "v0.2.0")
	}
}

func TestResolve_pseudoVersionFromVCS(t *testing.T) {
	prevOverride := Override
	prevRead := readBuildInfo
	Override = ""
	stamp := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{Version: "(devel)"},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "abcdef0123456789deadbeef"},
				{Key: "vcs.time", Value: stamp.Format(time.RFC3339)},
			},
		}, true
	}
	t.Cleanup(func() {
		Override = prevOverride
		readBuildInfo = prevRead
	})

	want := "v0.0.0-20260726120000-abcdef012345"
	if got := Resolve(); got != want {
		t.Fatalf("Resolve() = %q, want %q", got, want)
	}
}

func TestResolve_missing(t *testing.T) {
	prevOverride := Override
	prevRead := readBuildInfo
	Override = ""
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, true
	}
	t.Cleanup(func() {
		Override = prevOverride
		readBuildInfo = prevRead
	})

	if got := Resolve(); got != "" {
		t.Fatalf("Resolve() = %q, want empty", got)
	}
}
