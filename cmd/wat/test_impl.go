// wat test executes the user's compiled .wat/hooks binary (same cache path as
// wat run). The fixture report peeks hook_event_name from the payload; it does
// not replace hook execution.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/dialect"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

type testConfig struct {
	agent   string
	event   string
	fixture string
	verbose bool
}

type testDeps struct {
	runDeps
	readFixture func(path string, stdin io.Reader) ([]byte, error)
	writeReport io.Writer
}

type fixtureInfo struct {
	dialect string
	event   string
}

func defaultTestDeps() testDeps {
	rd := defaultRunDeps()
	return testDeps{
		runDeps: rd,
		readFixture: func(path string, stdin io.Reader) ([]byte, error) {
			if path == "-" {
				return io.ReadAll(stdin)
			}
			return rd.readFile(path)
		},
		writeReport: stdout,
	}
}

func runTest(cfg testConfig, deps testDeps) int {
	watDir, err := resolveWatDir(deps.runDeps)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "wat test: %v\n", err)
		return exitRuntimeFailure
	}

	payload, err := deps.readFixture(cfg.fixture, os.Stdin)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "wat test: read fixture: %v\n", err)
		return exitRuntimeFailure
	}
	if len(payload) == 0 {
		_, _ = fmt.Fprintln(stderr, "wat test: empty fixture")
		return exitRuntimeFailure
	}

	info, err := resolveFixture(cfg.agent, cfg.event, payload)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "wat test: %v\n", err)
		return exitRuntimeFailure
	}

	binPath, exitCode, ok := ensureHookBinary(watDir, deps.runDeps)
	if !ok {
		return exitCode
	}

	hookStdout, hookStderr, hookExit, err := execHookBinary(binPath, payload, cfg, deps.runDeps)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "wat test: exec %s: %v\n", binPath, err)
		return exitRuntimeFailure
	}

	writeTestReport(deps.writeReport, info, hookStdout, hookStderr, hookExit, cfg.verbose)
	return hookExit
}

func resolveFixture(agentFlag, _ string, payload []byte) (fixtureInfo, error) {
	agentDialect := dialect.Parse(agentFlag)
	if agentDialect == "" {
		return fixtureInfo{}, fmt.Errorf("unknown dialect (pass --agent)")
	}

	event, err := hookkit.PeekHookEventName(payload)
	if err != nil {
		return fixtureInfo{}, fmt.Errorf("decode: %w", err)
	}
	if event == "" {
		return fixtureInfo{}, fmt.Errorf("decode: hook_event_name is required")
	}
	if agentDialect == sdkcopilot.Dialect {
		if canonical, ok := sdkcopilot.ResolveCanonical(payload, event); ok {
			event = canonical
		}
	}
	switch agentDialect {
	case sdkclaude.Dialect, sdkcopilot.Dialect, sdkcursor.Dialect:
	default:
		return fixtureInfo{}, fmt.Errorf("unknown dialect (pass --agent)")
	}
	return fixtureInfo{dialect: agentDialect, event: event}, nil
}

func execHookBinary(binPath string, payload []byte, cfg testConfig, deps runDeps) (hookStdout, hookStderr []byte, exitCode int, err error) {
	cmd := deps.command(binPath)
	cmd.Stdin = bytes.NewReader(payload)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Env = append([]string(nil), os.Environ()...)
	if cfg.agent != "" {
		cmd.Env = append(cmd.Env, "WAT_AGENT="+cfg.agent)
	}

	runErr := deps.runCmd(cmd)
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			return outBuf.Bytes(), errBuf.Bytes(), exitErr.ExitCode(), nil
		}
		return nil, nil, 0, runErr
	}
	return outBuf.Bytes(), errBuf.Bytes(), exitOK, nil
}

func writeTestReport(w io.Writer, info fixtureInfo, hookStdout, hookStderr []byte, hookExit int, verbose bool) {
	_, _ = fmt.Fprintln(w, "fixture:")
	_, _ = fmt.Fprintf(w, "  agent: %s\n", info.dialect)
	_, _ = fmt.Fprintf(w, "  event: %s\n", info.event)

	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "hook:")
	if len(hookStdout) > 0 {
		_, _ = fmt.Fprintf(w, "  stdout: %s\n", strings.TrimSpace(string(hookStdout)))
	} else {
		_, _ = fmt.Fprintln(w, "  stdout: (empty)")
	}
	if decision := summarizeHookDecision(info.dialect, hookStdout); decision != "" {
		_, _ = fmt.Fprintf(w, "  decision: %s\n", decision)
	}
	if verbose && len(hookStderr) > 0 {
		_, _ = fmt.Fprintf(w, "  stderr: %s\n", strings.TrimSpace(string(hookStderr)))
	}
	_, _ = fmt.Fprintf(w, "  exit:   %d\n", hookExit)
}

func summarizeHookDecision(dialect string, hookStdout []byte) string {
	hookStdout = bytes.TrimSpace(hookStdout)
	if len(hookStdout) == 0 {
		return ""
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(hookStdout, &obj); err != nil {
		return ""
	}

	keys := decisionJSONKeys(dialect)
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

func decisionJSONKeys(dialect string) []string {
	switch dialect {
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
