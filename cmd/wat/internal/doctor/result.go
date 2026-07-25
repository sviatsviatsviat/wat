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

// PrintResult writes a single check result to w.
// Status labels are colored when w is a terminal and NO_COLOR is unset.
func PrintResult(w io.Writer, r Result) {
	label, color := statusLabel(r.Status)
	prefix := fmt.Sprintf("%-4s", label)
	if useColor(w) && color != "" {
		prefix = color + prefix + ansiReset
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

func useColor(w io.Writer) bool {
	if v, ok := os.LookupEnv("NO_COLOR"); ok && v != "" {
		return false
	}
	if v := os.Getenv("FORCE_COLOR"); v != "" && v != "0" {
		return true
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// PrintFailureSummary writes the doctor failure footer to w.
func PrintFailureSummary(w io.Writer, failCount int) {
	if failCount <= 0 {
		return
	}
	msg := fmt.Sprintf("wat doctor: %d check(s) failed", failCount)
	if useColor(w) {
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
