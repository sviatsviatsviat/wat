package doctor

import (
	"fmt"
	"go/version"
	"strings"
)

// ParseGoModDirective returns the go version from a go.mod file body.
func ParseGoModDirective(data []byte) (string, error) {
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go ") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "go "))
			if v == "" {
				return "", fmt.Errorf("empty go directive")
			}
			if i := strings.IndexByte(v, ' '); i >= 0 {
				v = v[:i]
			}
			return v, nil
		}
	}
	return "", fmt.Errorf("no go directive found")
}

// ParseInstalledGoVersion extracts the semver from go version output.
func ParseInstalledGoVersion(output string) (string, error) {
	fields := strings.Fields(strings.TrimSpace(output))
	for _, f := range fields {
		if strings.HasPrefix(f, "go") && len(f) > 2 {
			v := strings.TrimPrefix(f, "go")
			if v != "" {
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("no go version in output %q", strings.TrimSpace(output))
}

// GoVersionAtLeast reports whether installed satisfies required.
func GoVersionAtLeast(installed, required string) bool {
	return version.Compare(goVersionName(installed), goVersionName(required)) >= 0
}

func goVersionName(v string) string {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(v, "go") {
		return v
	}
	return "go" + v
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
