// wat test executes the user's compiled .wat/hooks binary (same cache path as
// wat run). Agent SDK decode plus inbound MapEvent build the fixture summary
// only; they do not replace hook execution.
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

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	agclaude "github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	agcopilot "github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	agcursor "github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
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

	ev, dialect, err := decodeFixtureSummary(cfg.agent, cfg.event, payload, deps.getenv)
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

	writeTestReport(deps.writeReport, ev, dialect, hookStdout, hookStderr, hookExit, cfg.verbose)
	return hookExit
}

func decodeFixtureSummary(agentFlag, eventHint string, payload []byte, getenv func(string) string) (*agnostic.Event, agnostic.Dialect, error) {
	dialect := agnostic.ParseDialect(agentFlag)
	if dialect == agnostic.Unknown {
		dialect = agnostic.Detect(payload, getenv)
	}
	if dialect == agnostic.Unknown {
		return nil, agnostic.Unknown, fmt.Errorf("unknown dialect (pass --agent or use a recognizable fixture)")
	}

	var ev *agnostic.Event
	var err error
	switch dialect {
	case agnostic.Claude:
		native, decErr := sdkclaude.Decode(payload)
		if decErr != nil {
			err = decErr
			break
		}
		ev = agclaude.MapEvent(native, payload)
	case agnostic.Copilot:
		native, decErr := sdkcopilot.Decode(payload, sdkcopilot.WithEvent(eventHint))
		if decErr != nil {
			err = decErr
			break
		}
		ev = agcopilot.MapEvent(native, payload)
	case agnostic.Cursor:
		native, decErr := sdkcursor.Decode(payload, sdkcursor.WithEvent(eventHint))
		if decErr != nil {
			err = decErr
			break
		}
		ev = agcursor.MapEvent(native, payload)
	default:
		return nil, dialect, fmt.Errorf("unknown dialect (pass --agent or use a recognizable fixture)")
	}
	if err != nil {
		if dialect == agnostic.Copilot && eventHint == "" {
			return nil, dialect, fmt.Errorf("decode: %w (Copilot camelCase payloads require --event)", err)
		}
		return nil, dialect, fmt.Errorf("decode: %w", err)
	}
	return ev, dialect, nil
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
	if cfg.event != "" {
		cmd.Env = append(cmd.Env, "WAT_EVENT="+cfg.event)
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

func writeTestReport(w io.Writer, ev *agnostic.Event, dialect agnostic.Dialect, hookStdout, hookStderr []byte, hookExit int, verbose bool) {
	_, _ = fmt.Fprintln(w, "event:")
	_, _ = fmt.Fprintf(w, "  agent: %s\n", dialect)
	_, _ = fmt.Fprintf(w, "  kind:  %s\n", ev.Kind)
	_, _ = fmt.Fprintf(w, "  name:  %s\n", ev.Name)
	if ev.Tool != nil {
		_, _ = fmt.Fprintf(w, "  tool:  %s\n", ev.Tool.Name)
		if ev.Tool.Shell != "" {
			_, _ = fmt.Fprintf(w, "  shell: %s\n", ev.Tool.Shell)
		}
		if verbose {
			if ev.Tool.Native != "" && ev.Tool.Native != ev.Tool.Name {
				_, _ = fmt.Fprintf(w, "  tool_native: %s\n", ev.Tool.Native)
			}
		}
	}
	if verbose {
		if ev.Session != "" {
			_, _ = fmt.Fprintf(w, "  session: %s\n", ev.Session)
		}
		if ev.Cwd != "" {
			_, _ = fmt.Fprintf(w, "  cwd:     %s\n", ev.Cwd)
		}
	}

	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "hook:")
	if len(hookStdout) > 0 {
		_, _ = fmt.Fprintf(w, "  stdout: %s\n", strings.TrimSpace(string(hookStdout)))
	} else {
		_, _ = fmt.Fprintln(w, "  stdout: (empty)")
	}
	if decision := summarizeHookDecision(dialect, hookStdout); decision != "" {
		_, _ = fmt.Fprintf(w, "  decision: %s\n", decision)
	}
	if verbose && len(hookStderr) > 0 {
		_, _ = fmt.Fprintf(w, "  stderr: %s\n", strings.TrimSpace(string(hookStderr)))
	}
	_, _ = fmt.Fprintf(w, "  exit:   %d\n", hookExit)
}

func summarizeHookDecision(dialect agnostic.Dialect, hookStdout []byte) string {
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

func decisionJSONKeys(dialect agnostic.Dialect) []string {
	switch dialect {
	case agnostic.Claude:
		return []string{"permissionDecision", "decision", "permission"}
	case agnostic.Copilot:
		return []string{"permissionDecision", "decision", "permission"}
	case agnostic.Cursor:
		return []string{"permission", "permissionDecision", "decision"}
	default:
		return []string{"permissionDecision", "permission", "decision"}
	}
}
