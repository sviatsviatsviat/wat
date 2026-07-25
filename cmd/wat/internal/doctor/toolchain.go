package doctor

import (
	"fmt"
	"path/filepath"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/project"
)

// ToolchainGoOnPath verifies go is available on PATH.
func ToolchainGoOnPath(deps Deps) []Result {
	if _, err := deps.LookPath("go"); err != nil {
		return []Result{{
			Group:   "toolchain",
			Status:  Fail,
			Message: "go not found on PATH",
			Fix:     "install Go and ensure go is on PATH",
		}}
	}
	return []Result{{
		Group:   "toolchain",
		Status:  Pass,
		Message: "go on PATH",
	}}
}

// ToolchainGoVersion verifies the installed Go version satisfies .wat/go.mod.
func ToolchainGoVersion(deps Deps, ctx Context) []Result {
	if ctx.WatErr != nil {
		return []Result{{
			Group:   "toolchain",
			Status:  Fail,
			Message: "no .wat/go.mod found",
			Fix:     "run wat init",
		}}
	}
	goModPath := filepath.Join(ctx.WatDir, project.GoModFile)
	data, err := deps.ReadFile(goModPath)
	if err != nil {
		return []Result{{
			Group:   "toolchain",
			Status:  Fail,
			Message: fmt.Sprintf("read %s failed", goModPath),
			Fix:     "run wat init",
		}}
	}
	required, err := ParseGoModDirective(data)
	if err != nil {
		return []Result{{
			Group:   "toolchain",
			Status:  Fail,
			Message: fmt.Sprintf("parse go directive in %s: %v", goModPath, err),
			Fix:     "fix the go directive in .wat/go.mod",
		}}
	}

	cmd := deps.Command("go", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []Result{{
			Group:   "toolchain",
			Status:  Fail,
			Message: "go version failed",
			Fix:     "install Go and ensure go is on PATH",
		}}
	}
	installed, err := ParseInstalledGoVersion(string(out))
	if err != nil {
		return []Result{{
			Group:   "toolchain",
			Status:  Fail,
			Message: fmt.Sprintf("parse go version output: %v", err),
			Fix:     "install a supported Go toolchain",
		}}
	}
	if !GoVersionAtLeast(installed, required) {
		return []Result{{
			Group:   "toolchain",
			Status:  Fail,
			Message: fmt.Sprintf("installed %s does not satisfy go %s in .wat/go.mod", installed, required),
			Fix:     fmt.Sprintf("upgrade Go to >= %s", required),
		}}
	}
	return []Result{{
		Group:   "toolchain",
		Status:  Pass,
		Message: fmt.Sprintf("%s satisfies go %s directive", installed, required),
	}}
}
