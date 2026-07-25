package doctor

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/buildcache"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookmanifest"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func authoredManifest(deps Deps, ctx Context) (run.Manifest, []Result) {
	if ctx.WatErr != nil {
		return run.Manifest{}, []Result{{
			Group:   "script",
			Status:  Fail,
			Message: "cannot inspect registrations without .wat/ project",
			Fix:     "run wat init",
		}}
	}

	bc := buildcache.Adapt(deps.Getenv, deps.Stat, deps.ReadDir, deps.ReadFile, deps.MkdirAll, deps.WriteFile, deps.Command)
	var diagnostics bytes.Buffer
	loadManifest := deps.LoadManifest
	if loadManifest == nil {
		loadManifest = hookmanifest.Load
	}
	manifest, err := loadManifest(ctx.WatDir, deps.WatVersion, bc, &diagnostics)
	if err != nil {
		message := fmt.Sprintf("load authored hook registrations: %v", err)
		if detail := strings.TrimSpace(diagnostics.String()); detail != "" {
			message += " (" + firstLine(detail) + ")"
		}
		return run.Manifest{}, []Result{{
			Group:   "script",
			Status:  Fail,
			Message: message,
			Fix:     "fix .wat/hooks.go and re-run wat doctor",
		}}
	}

	return manifest, []Result{{
		Group:   "script",
		Status:  Pass,
		Message: fmt.Sprintf("%d native hook registration(s) loaded", len(manifest.Registrations)),
	}}
}
