package doctor

// Run executes all doctor checks in order and returns the combined results.
func Run(deps Deps, ctx Context) []Result {
	var results []Result
	results = append(results, ToolchainGoOnPath(deps)...)
	results = append(results, ToolchainGoVersion(deps, ctx)...)
	results = append(results, ScriptFiles(deps, ctx)...)
	results = append(results, ScriptBuild(deps, ctx)...)
	results = append(results, CacheWritable(deps, ctx)...)
	results = append(results, CacheWarm(deps, ctx)...)
	results = append(results, Install(deps, ctx)...)
	return results
}
