package doctor

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/project"
)

// scriptBuildTimeout bounds the doctor go build probe.
const scriptBuildTimeout = 2 * time.Minute

// ScriptFiles verifies .wat/hooks.go and go.mod are present.
func ScriptFiles(deps Deps, ctx Context) []Result {
	if ctx.WatErr != nil {
		return []Result{{
			Group:   "script",
			Status:  Fail,
			Message: "no .wat/ project found",
			Fix:     "run wat init",
		}}
	}
	proj := project.Deps{Getenv: deps.Getenv, Getwd: deps.Getwd, Stat: deps.Stat}
	if err := project.MustHaveFiles(ctx.WatDir, proj); err != nil {
		return []Result{{
			Group:   "script",
			Status:  Fail,
			Message: ".wat/hooks.go or .wat/go.mod missing",
			Fix:     "run wat init",
		}}
	}
	return []Result{{
		Group:   "script",
		Status:  Pass,
		Message: ".wat/hooks.go and go.mod present",
	}}
}

// ScriptBuild verifies go build succeeds in .wat/.
func ScriptBuild(deps Deps, ctx Context) []Result {
	if ctx.WatErr != nil {
		return []Result{{
			Group:   "script",
			Status:  Fail,
			Message: "cannot compile without .wat/ project",
			Fix:     "run wat init",
		}}
	}
	probeDir := filepath.Join(ctx.WatDir, ".cache", ".doctor-probe")
	binName := "hooks"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(probeDir, binName)
	if err := deps.MkdirAll(probeDir, 0o755); err != nil {
		return []Result{{
			Group:   "script",
			Status:  Fail,
			Message: fmt.Sprintf("create %s failed", probeDir),
			Fix:     "fix permissions on .wat/.cache/",
		}}
	}
	defer func() { _ = deps.Remove(binPath); _ = deps.Remove(probeDir) }()

	cmd := deps.Command("go", "build", "-o", binPath)
	cmd.Dir = ctx.WatDir
	out, err := combinedOutputWithTimeout(cmd, scriptBuildTimeout)
	if err == nil {
		return []Result{{
			Group:   "script",
			Status:  Pass,
			Message: "go build in .wat/ succeeds",
		}}
	}
	if errors.Is(err, errCommandTimeout) {
		return []Result{{
			Group:   "script",
			Status:  Fail,
			Message: fmt.Sprintf("go build timed out after %s", scriptBuildTimeout),
			Fix:     "fix hang or slow compile in .wat/hooks.go",
		}}
	}
	msg := "go build failed in .wat/"
	if len(out) > 0 {
		trimmed := strings.TrimSpace(string(out))
		if trimmed != "" {
			msg = msg + " (" + firstLine(trimmed) + ")"
		}
	}
	return []Result{{
		Group:   "script",
		Status:  Fail,
		Message: msg,
		Fix:     "fix compile errors in .wat/hooks.go",
	}}
}

var errCommandTimeout = errors.New("command timed out")

func combinedOutputWithTimeout(cmd *exec.Cmd, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type result struct {
		out []byte
		err error
	}
	done := make(chan result, 1)
	go func() {
		out, err := cmd.CombinedOutput()
		done <- result{out: out, err: err}
	}()

	select {
	case r := <-done:
		return r.out, r.err
	case <-ctx.Done():
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		r := <-done
		if r.err != nil && len(r.out) > 0 {
			return r.out, fmt.Errorf("%w: %v", errCommandTimeout, r.err)
		}
		return r.out, errCommandTimeout
	}
}
