package doctor

import (
	"fmt"
	"io"
	"os"
)

// Status is the outcome of a single doctor check.
type Status int

const (
	// Pass means the check succeeded.
	Pass Status = iota
	// Fail means the check failed and should fail the doctor run.
	Fail
	// Warn means the check found a non-fatal issue.
	Warn
)

// Result is one doctor check outcome with an optional fix suggestion.
type Result struct {
	Group   string
	Status  Status
	Message string
	Fix     string
}

const (
	ansiReset  = "\033[0m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiRed    = "\033[31m"
)

// ColorEnabled reports whether doctor output should use ANSI colors for w.
// It honors NO_COLOR and FORCE_COLOR via deps, then falls back to deps.IsTerminal.
func ColorEnabled(deps Deps, w io.Writer) bool {
	lookup := deps.LookupEnv
	if lookup == nil {
		lookup = os.LookupEnv
	}
	if v, ok := lookup("NO_COLOR"); ok && v != "" {
		return false
	}
	getenv := deps.Getenv
	if getenv == nil {
		getenv = os.Getenv
	}
	if v := getenv("FORCE_COLOR"); v != "" && v != "0" {
		return true
	}
	if deps.IsTerminal != nil {
		return deps.IsTerminal(w)
	}
	return false
}

// PrintResult writes a single check result to w.
// When color is true, status labels use ANSI colors.
func PrintResult(w io.Writer, r Result, color bool) {
	label, ansi := statusLabel(r.Status)
	prefix := fmt.Sprintf("%-4s", label)
	if color && ansi != "" {
		prefix = ansi + prefix + ansiReset
	}
	_, _ = fmt.Fprintf(w, "%s  %-9s  %s\n", prefix, r.Group, r.Message)
	if r.Fix != "" {
		_, _ = fmt.Fprintf(w, "      fix: %s\n", r.Fix)
	}
}

func statusLabel(s Status) (label, color string) {
	switch s {
	case Pass:
		return "PASS", ansiGreen
	case Fail:
		return "FAIL", ansiRed
	case Warn:
		return "WARN", ansiYellow
	default:
		return "????", ""
	}
}

// PrintFailureSummary writes the doctor failure footer to w.
func PrintFailureSummary(w io.Writer, failCount int, color bool) {
	if failCount <= 0 {
		return
	}
	msg := fmt.Sprintf("wat doctor: %d check(s) failed", failCount)
	if color {
		msg = ansiRed + msg + ansiReset
	}
	_, _ = fmt.Fprintf(w, "\n%s\n", msg)
}

// FailCount returns how many results have Status Fail.
func FailCount(results []Result) int {
	n := 0
	for _, r := range results {
		if r.Status == Fail {
			n++
		}
	}
	return n
}
