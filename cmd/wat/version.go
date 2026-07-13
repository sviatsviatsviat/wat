package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

func watModuleVersion() string {
	info, ok := debug.ReadBuildInfo()
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

var watModuleVersionFn = watModuleVersion
