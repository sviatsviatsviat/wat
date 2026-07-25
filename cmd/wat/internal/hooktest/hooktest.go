package hooktest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/buildcache"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/dialect"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookmanifest"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/project"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// Exit codes returned by Run for non-hook failures.
const (
	ExitOK             = 0
	ExitBuildFailed    = 1
	ExitRuntimeFailure = 3
)

// Config holds options for Run.
type Config struct {
	// Agent is the required agent dialect for the fixture report.
	Agent string
	// Fixture is a path to fixture JSON, or "-" for stdin.
	Fixture string
	// Expect is an optional path to an expect JSON document. When empty, Run
	// loads the conventional sidecar (<fixture>.expect.json) if it exists.
	Expect string
	// Verbose prints hook stderr in the report.
	Verbose bool
}

// Deps holds injectable dependencies for Run.
type Deps struct {
	Getenv      func(string) string
	Getwd       func() (string, error)
	Stat        func(string) (os.FileInfo, error)
	ReadDir     func(string) ([]os.DirEntry, error)
	ReadFile    func(string) ([]byte, error)
	MkdirAll    func(string, os.FileMode) error
	WriteFile   func(string, []byte, os.FileMode) error
	Command     func(string, ...string) *exec.Cmd
	RunCmd      func(*exec.Cmd) error
	ReadFixture func(path string, stdin io.Reader) ([]byte, error)
	WriteReport io.Writer
}

// DefaultDeps returns production dependencies. writeReport receives the fixture report.
// ReadFixture is left nil so Run reads fixture files through Deps.ReadFile; tests may
// set ReadFixture to override that behavior.
func DefaultDeps(writeReport io.Writer) Deps {
	return Deps{
		Getenv:    os.Getenv,
		Getwd:     os.Getwd,
		Stat:      os.Stat,
		ReadDir:   os.ReadDir,
		ReadFile:  os.ReadFile,
		MkdirAll:  os.MkdirAll,
		WriteFile: os.WriteFile,
		Command:   exec.Command,
		RunCmd: func(cmd *exec.Cmd) error {
			return cmd.Run()
		},
		WriteReport: writeReport,
	}
}

// FixtureInfo is the agent/event summary for a fixture report.
type FixtureInfo struct {
	// Dialect is the normalized agent dialect.
	Dialect string
	// Event is the hook_event_name from the fixture payload.
	Event string
}

// Run executes the hooks binary against a fixture and writes a report.
func Run(cfg Config, version string, deps Deps, stdin io.Reader, errOut io.Writer) int {
	watDir, err := project.Resolve(project.Deps{
		Getenv: deps.Getenv,
		Getwd:  deps.Getwd,
		Stat:   deps.Stat,
	})
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat test: %v\n", err)
		return ExitRuntimeFailure
	}

	payload, err := loadFixture(deps, cfg.Fixture, stdin)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat test: read fixture: %v\n", err)
		return ExitRuntimeFailure
	}
	if len(payload) == 0 {
		_, _ = fmt.Fprintln(errOut, "wat test: empty fixture")
		return ExitRuntimeFailure
	}

	info, err := ResolveFixture(cfg.Agent, payload)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat test: %v\n", err)
		return ExitRuntimeFailure
	}

	bc := buildcache.Adapt(deps.Getenv, deps.Stat, deps.ReadDir, deps.ReadFile, deps.MkdirAll, deps.WriteFile, deps.Command)
	binPath, manifest, err := hookmanifest.EnsureAndLoad(watDir, version, bc, errOut)
	if err != nil {
		if errors.Is(err, buildcache.ErrBuildFailed) {
			return ExitBuildFailed
		}
		_, _ = fmt.Fprintf(errOut, "wat test: %v\n", err)
		return ExitRuntimeFailure
	}
	if !manifest.Has(info.Dialect, info.Event) {
		_, _ = fmt.Fprintf(errOut, "wat test: no %s %s handler is registered\n", info.Dialect, info.Event)
		return ExitRuntimeFailure
	}

	hookStdout, hookStderr, hookExit, err := execHookBinary(binPath, payload, deps)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat test: exec %s: %v\n", binPath, err)
		return ExitRuntimeFailure
	}

	WriteReport(deps.WriteReport, info, hookStdout, hookStderr, hookExit, cfg.Verbose)

	expectPath, err := ResolveExpectPath(cfg.Fixture, cfg.Expect, deps.Stat)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat test: %v\n", err)
		return ExitRuntimeFailure
	}
	if expectPath == "" {
		return hookExit
	}

	exp, err := LoadExpect(expectPath, deps.ReadFile)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat test: %v\n", err)
		return ExitRuntimeFailure
	}
	fails := CheckExpect(exp, info.Dialect, hookStdout, hookExit)
	writeExpectReport(deps.WriteReport, expectPath, fails)
	if len(fails) > 0 {
		for _, fail := range fails {
			_, _ = fmt.Fprintf(errOut, "wat test: expect failed: %s\n", fail)
		}
		return ExitExpectFailed
	}
	return ExitOK
}

// ResolveFixture peeks hook_event_name and validates the agent dialect.
func ResolveFixture(agentFlag string, payload []byte) (FixtureInfo, error) {
	agentDialect := dialect.Parse(agentFlag)
	if agentDialect == "" {
		return FixtureInfo{}, fmt.Errorf("unknown dialect (pass --agent)")
	}

	event, err := hookkit.PeekHookEventName(payload)
	if err != nil {
		return FixtureInfo{}, fmt.Errorf("decode: %w", err)
	}
	if event == "" {
		return FixtureInfo{}, fmt.Errorf("decode: hook_event_name is required")
	}
	switch agentDialect {
	case sdkclaude.Dialect, sdkcopilot.Dialect, sdkcursor.Dialect:
	default:
		return FixtureInfo{}, fmt.Errorf("unknown dialect (pass --agent)")
	}
	return FixtureInfo{Dialect: agentDialect, Event: event}, nil
}

func loadFixture(deps Deps, path string, stdin io.Reader) ([]byte, error) {
	if deps.ReadFixture != nil {
		return deps.ReadFixture(path, stdin)
	}
	if path == "-" {
		return io.ReadAll(stdin)
	}
	return deps.ReadFile(path)
}

func execHookBinary(binPath string, payload []byte, deps Deps) (hookStdout, hookStderr []byte, exitCode int, err error) {
	cmd := deps.Command(binPath)
	cmd.Stdin = bytes.NewReader(payload)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	runErr := deps.RunCmd(cmd)
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			return outBuf.Bytes(), errBuf.Bytes(), exitErr.ExitCode(), nil
		}
		return nil, nil, 0, runErr
	}
	return outBuf.Bytes(), errBuf.Bytes(), ExitOK, nil
}

// WriteReport writes the fixture/hook summary report.
func WriteReport(w io.Writer, info FixtureInfo, hookStdout, hookStderr []byte, hookExit int, verbose bool) {
	_, _ = fmt.Fprintln(w, "fixture:")
	_, _ = fmt.Fprintf(w, "  agent: %s\n", info.Dialect)
	_, _ = fmt.Fprintf(w, "  event: %s\n", info.Event)

	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "hook:")
	if len(hookStdout) > 0 {
		_, _ = fmt.Fprintf(w, "  stdout: %s\n", strings.TrimSpace(string(hookStdout)))
	} else {
		_, _ = fmt.Fprintln(w, "  stdout: (empty)")
	}
	if decision := summarizeHookDecision(info.Dialect, hookStdout); decision != "" {
		_, _ = fmt.Fprintf(w, "  decision: %s\n", decision)
	}
	if verbose && len(hookStderr) > 0 {
		_, _ = fmt.Fprintf(w, "  stderr: %s\n", strings.TrimSpace(string(hookStderr)))
	}
	_, _ = fmt.Fprintf(w, "  exit:   %d\n", hookExit)
}

func writeExpectReport(w io.Writer, path string, fails []string) {
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "expect:")
	_, _ = fmt.Fprintf(w, "  file:   %s\n", path)
	if len(fails) == 0 {
		_, _ = fmt.Fprintln(w, "  status: pass")
		return
	}
	_, _ = fmt.Fprintln(w, "  status: fail")
	for _, fail := range fails {
		_, _ = fmt.Fprintf(w, "  - %s\n", fail)
	}
}

func summarizeHookDecision(dialectName string, hookStdout []byte) string {
	hookStdout = bytes.TrimSpace(hookStdout)
	if len(hookStdout) == 0 {
		return ""
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(hookStdout, &obj); err != nil {
		return ""
	}

	keys := decisionJSONKeys(dialectName)
	for _, key := range keys {
		raw, ok := obj[key]
		if !ok {
			continue
		}
		var s string
		if err := json.Unmarshal(raw, &s); err == nil && s != "" {
			return s
		}
	}
	return ""
}

func decisionJSONKeys(dialectName string) []string {
	switch dialectName {
	case sdkclaude.Dialect:
		return []string{"permissionDecision", "decision", "permission"}
	case sdkcopilot.Dialect:
		return []string{"permission_decision", "decision", "permission"}
	case sdkcursor.Dialect:
		return []string{"permission", "permission_decision", "decision"}
	default:
		return []string{"permissionDecision", "permission_decision", "permission", "decision"}
	}
}
