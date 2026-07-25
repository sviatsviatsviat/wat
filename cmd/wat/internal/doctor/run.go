package doctor

import "fmt"

// Run executes all doctor checks in order and returns the combined results.
func Run(deps Deps, ctx Context) []Result {
	var results []Result
	results = append(results, ToolchainGoOnPath(deps)...)
	results = append(results, ToolchainGoVersion(deps, ctx)...)
	results = append(results, ScriptFiles(deps, ctx)...)
	results = append(results, ScriptBuild(deps, ctx)...)
	results = append(results, CacheWritable(deps, ctx)...)
	results = append(results, CacheWarm(deps, ctx)...)
	manifest, manifestResults := authoredManifest(deps, ctx)
	results = append(results, manifestResults...)
	ctx.Manifest = manifest
	if FailCount(manifestResults) > 0 {
		ctx.ManifestErr = fmt.Errorf("authored hook manifest unavailable")
	}
	results = append(results, Install(deps, ctx)...)
	return results
}
