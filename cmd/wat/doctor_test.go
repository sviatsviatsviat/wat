package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/buildcache"
	hostclaude "github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/claude"
	hostcursor "github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/installcfg"
	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
	sdkrun "github.com/sviatsviatsviat/wat/sdk/run"
)

func TestRunDoctor_missingGo(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{buildOK: true})
	deps.LookPath = func(name string) (string, error) {
		if name == "go" {
			return "", os.ErrNotExist
		}
		return exec.LookPath(name)
	}

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitCheckFailed {
		t.Fatalf("exit = %d, want %d\n%s", code, exitCheckFailed, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "go not found on PATH") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_goVersionTooOld(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.25.0 linux/amd64",
		buildOK:   true,
	})

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitCheckFailed {
		t.Fatalf("exit = %d, want %d\n%s", code, exitCheckFailed, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "does not satisfy go 1.26") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_missingWatProject(t *testing.T) {
	dir := t.TempDir()
	deps := doctorTestDeps(t, dir, filepath.Join(dir, "wat"), doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   true,
	})

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitCheckFailed {
		t.Fatalf("exit = %d, want %d\n%s", code, exitCheckFailed, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "no .wat/ project found") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_compileError(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	if err := os.WriteFile(filepath.Join(project, ".wat", "hooks.go"), []byte("package main\n\nfunc main() {\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   false,
		buildErr:  "syntax error",
	})

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitCheckFailed {
		t.Fatalf("exit = %d, want %d\n%s", code, exitCheckFailed, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "go build failed") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_cacheNotWritable(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   true,
	})
	cacheDir := filepath.Join(project, ".wat", ".cache")
	deps.WriteFile = func(path string, data []byte, perm os.FileMode) error {
		if strings.HasPrefix(path, cacheDir) {
			return os.ErrPermission
		}
		return os.WriteFile(path, data, perm)
	}

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitCheckFailed {
		t.Fatalf("exit = %d, want %d\n%s", code, exitCheckFailed, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), ".wat/.cache/") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_coldCache(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   true,
	})

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitOK {
		t.Fatalf("exit = %d, want %d\n%s", code, exitOK, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "no cached binary") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_watNotOnPath(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   true,
	})
	deps.LookPath = func(name string) (string, error) {
		if name == "wat" {
			return "", os.ErrNotExist
		}
		return exec.LookPath(name)
	}

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitOK {
		t.Fatalf("exit = %d, want %d\n%s", code, exitOK, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "wat not found on PATH") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_missingInstallEntry(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   true,
	})
	cursorPath := filepath.Join(project, ".cursor", "hooks.json")
	var f hostcursor.File
	data, err := os.ReadFile(cursorPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	delete(f.Hooks, cursor.EventSessionEnd)
	b, _ := json.MarshalIndent(f, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(cursorPath, b, 0o644); err != nil {
		t.Fatal(err)
	}

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitOK {
		t.Fatalf("exit = %d, want %d\n%s", code, exitOK, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "missing hook entry") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_invalidEventInConfig(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	claudePath := filepath.Join(project, ".claude", "settings.json")
	event := claude.EventPreToolUse
	badCmd := watAbs + " run --agent claude --event NotARealEvent"
	settings := hostclaude.Settings{
		Hooks: map[string][]hostclaude.MatcherGroup{
			event: {{
				Hooks: []json.RawMessage{
					mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: badCmd}),
				},
			}},
		},
	}
	b, _ := json.MarshalIndent(settings, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(claudePath, b, 0o644); err != nil {
		t.Fatal(err)
	}

	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   true,
	})

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitCheckFailed {
		t.Fatalf("exit = %d, want %d\n%s", code, exitCheckFailed, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "unregistered --event") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_disableAllHooks(t *testing.T) {
	project, watAbs := doctorTestProject(t)
	claudePath := filepath.Join(project, ".claude", "settings.json")
	var settings hostclaude.Settings
	data, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatal(err)
	}
	settings.DisableAllHooks = true
	b, _ := json.MarshalIndent(settings, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(claudePath, b, 0o644); err != nil {
		t.Fatal(err)
	}

	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   true,
	})
	warmDoctorCache(t, project, deps)

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitOK {
		t.Fatalf("exit = %d, want %d\n%s", code, exitOK, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "disableAllHooks is true") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

func TestRunDoctor_noInstallConfigs(t *testing.T) {
	project, watAbs := doctorTestProject(t, doctorProjectOpts{skipInstall: true})
	deps := doctorTestDeps(t, project, watAbs, doctorTestGoDeps{
		goVersion: "go version go1.26.0 linux/amd64",
		buildOK:   true,
	})

	outBuf := captureStdout(t)
	code := runDoctor(deps)
	if code != exitOK {
		t.Fatalf("exit = %d, want %d\n%s", code, exitOK, outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "hook config file missing") {
		t.Fatalf("output: %s", outBuf.String())
	}
}

type doctorTestGoDeps struct {
	goVersion string
	buildOK   bool
	buildErr  string
}

type doctorProjectOpts struct {
	skipInstall bool
}

func doctorTestProject(t *testing.T, opts ...doctorProjectOpts) (project, watAbs string) {
	t.Helper()
	var o doctorProjectOpts
	if len(opts) > 0 {
		o = opts[0]
	}

	project = t.TempDir()
	watDir := filepath.Join(project, ".wat")
	if err := os.MkdirAll(filepath.Join(watDir, ".cache"), 0o755); err != nil {
		t.Fatal(err)
	}
	hooksSource := `package hooks

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

var Hooks = []run.Hooks{
	agnostic.UseHooks().OnSessionEnd(func(context.Context, agnostic.SessionEndEvent) error {
		return nil
	}),
}
`
	if err := os.WriteFile(filepath.Join(watDir, "hooks.go"), []byte(hooksSource), 0o644); err != nil {
		t.Fatal(err)
	}
	goMod := "module wat-hooks\n\ngo 1.26\n\nrequire github.com/sviatsviatsviat/wat v0.0.0\n"
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	watAbs = filepath.Join(project, "bin", "wat")
	if !o.skipInstall {
		deps := installcfg.DefaultDeps()
		deps.Getwd = func() (string, error) { return project, nil }
		if err := installcfg.Install(installcfg.Config{
			Agents:   installcfg.AgentPlan{Claude: true, Copilot: true, Cursor: true},
			WatPath:  watAbs,
			Manifest: doctorTestManifest(),
		}, deps); err != nil {
			t.Fatalf("installProject: %v", err)
		}
	}
	return project, watAbs
}

func doctorTestDeps(t *testing.T, project, watAbs string, goCfg doctorTestGoDeps) doctorDeps {
	t.Helper()
	deps := defaultDoctorDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	deps.LookPath = func(name string) (string, error) {
		switch name {
		case "wat":
			return watAbs, nil
		case "go":
			return "go", nil
		default:
			return exec.LookPath(name)
		}
	}
	deps.Command = func(name string, args ...string) *exec.Cmd {
		if name == "go" && len(args) >= 2 && args[0] == "env" && args[1] == "GOVERSION" {
			return exec.Command("echo", "go1.26.0")
		}
		if name == "go" && len(args) > 0 && args[0] == "version" {
			return exec.Command("echo", goCfg.goVersion)
		}
		if name == "go" && len(args) >= 3 && args[0] == "build" {
			if goCfg.buildOK {
				out := args[2]
				return exec.Command("sh", "-c", "mkdir -p "+filepath.Dir(out)+" && touch "+out)
			}
			msg := goCfg.buildErr
			if msg == "" {
				msg = "build failed"
			}
			return exec.Command("sh", "-c", "echo "+msg+" >&2; exit 1")
		}
		return exec.Command(name, args...)
	}
	deps.LoadManifest = func(string, string, buildcache.Deps, io.Writer) (sdkrun.Manifest, error) {
		return *doctorTestManifest(), nil
	}
	return deps
}

func doctorTestManifest() *sdkrun.Manifest {
	return &sdkrun.Manifest{
		Version: 1,
		Registrations: []sdkrun.Registration{
			{Dialect: claude.Dialect, Event: claude.EventSessionEnd, HandlerCount: 1},
			{Dialect: copilot.Dialect, Event: copilot.EventSessionEnd, HandlerCount: 1},
			{Dialect: cursor.Dialect, Event: cursor.EventSessionEnd, HandlerCount: 1},
		},
	}
}

func warmDoctorCache(t *testing.T, project string, deps doctorDeps) {
	t.Helper()
	watDir := filepath.Join(project, ".wat")
	bc := buildcache.Adapt(deps.Getenv, deps.Stat, deps.ReadDir, deps.ReadFile, deps.MkdirAll, deps.WriteFile, deps.Command)
	key, err := buildcache.CacheKey(watDir, deps.WatVersion, bc)
	if err != nil {
		t.Fatal(err)
	}
	binPath := buildcache.BinaryPath(watDir, key)
	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binPath, []byte("binary"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func captureStdout(t *testing.T) *bytes.Buffer {
	t.Helper()
	prevStdout := stdout
	buf := &bytes.Buffer{}
	stdout = buf
	t.Cleanup(func() { stdout = prevStdout })
	return buf
}
