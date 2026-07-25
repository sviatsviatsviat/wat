package doctor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/buildcache"
)

// cacheWarmFix tells users how to populate the hooks build cache.
const cacheWarmFix = "run wat test with a project fixture, or otherwise warm the cache for this project"

// CacheWritable verifies .wat/.cache/ exists and is writable.
func CacheWritable(deps Deps, ctx Context) []Result {
	if ctx.WatErr != nil {
		return []Result{{
			Group:   "cache",
			Status:  Fail,
			Message: "cannot check cache without .wat/ project",
			Fix:     "run wat init",
		}}
	}
	cacheDir := filepath.Join(ctx.WatDir, ".cache")
	if err := deps.MkdirAll(cacheDir, 0o755); err != nil {
		return []Result{{
			Group:   "cache",
			Status:  Fail,
			Message: fmt.Sprintf("create %s failed", cacheDir),
			Fix:     "fix permissions on .wat/.cache/",
		}}
	}
	probe := filepath.Join(cacheDir, ".doctor-write-probe")
	if err := deps.WriteFile(probe, []byte("ok"), 0o644); err != nil {
		return []Result{{
			Group:   "cache",
			Status:  Fail,
			Message: fmt.Sprintf("write probe in %s failed", cacheDir),
			Fix:     "fix permissions on .wat/.cache/",
		}}
	}
	_ = deps.Remove(probe)
	return []Result{{
		Group:   "cache",
		Status:  Pass,
		Message: ".wat/.cache/ writable",
	}}
}

// CacheWarm warns when no cached binary exists for current hook sources.
func CacheWarm(deps Deps, ctx Context) []Result {
	if ctx.WatErr != nil {
		return nil
	}
	bc := buildcache.Adapt(deps.Getenv, deps.Stat, deps.ReadDir, deps.ReadFile, deps.MkdirAll, deps.Command)
	key, err := buildcache.CacheKey(ctx.WatDir, deps.WatVersion, bc)
	if err != nil {
		return []Result{{
			Group:   "cache",
			Status:  Warn,
			Message: "cannot compute cache key for current hook sources",
			Fix:     cacheWarmFix,
		}}
	}
	binPath := buildcache.BinaryPath(ctx.WatDir, key)
	if _, err := deps.Stat(binPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Result{{
				Group:   "cache",
				Status:  Warn,
				Message: "no cached binary for current hook sources",
				Fix:     cacheWarmFix,
			}}
		}
		return []Result{{
			Group:   "cache",
			Status:  Fail,
			Message: fmt.Sprintf("stat %s failed", binPath),
			Fix:     "fix permissions on .wat/.cache/",
		}}
	}
	return nil
}
