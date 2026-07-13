package checks

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

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
	if deps.MustHaveWatFiles == nil {
		return []Result{{
			Group:   "script",
			Status:  Fail,
			Message: "internal error: MustHaveWatFiles not configured",
			Fix:     "report a bug in wat",
		}}
	}
	if err := deps.MustHaveWatFiles(ctx.WatDir); err != nil {
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
	out, err := cmd.CombinedOutput()
	if err == nil {
		return []Result{{
			Group:   "script",
			Status:  Pass,
			Message: "go build in .wat/ succeeds",
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
