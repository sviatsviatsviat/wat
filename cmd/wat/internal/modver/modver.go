package modver

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

// Override replaces build-info module version resolution when non-empty.
// E2E sets it via -ldflags so linked git worktrees (which often omit VCS
// stamps) still exercise init/cache/version with a known pin.
var Override string

var readBuildInfo = debug.ReadBuildInfo

// Resolve returns the module version string used for wat init pinning and the
// hook build cache key, or "" when neither a module version nor VCS stamp is
// available.
func Resolve() string {
	if v := strings.TrimSpace(Override); v != "" {
		return v
	}
	info, ok := readBuildInfo()
	if !ok {
		return ""
	}
	if v := strings.TrimSpace(info.Main.Version); v != "" && v != "(devel)" {
		return v
	}

	var revision string
	var t time.Time
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			parsed, err := time.Parse(time.RFC3339, s.Value)
			if err == nil {
				t = parsed.UTC()
			}
		}
	}
	if revision == "" || t.IsZero() {
		return ""
	}
	short := revision
	if len(short) > 12 {
		short = short[:12]
	}
	return fmt.Sprintf("v0.0.0-%s-%s", t.Format("20060102150405"), short)
}
