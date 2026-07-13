package checks

import (
	"fmt"
	"io"
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

// PrintResult writes a single check result to w.
func PrintResult(w io.Writer, r Result) {
	var prefix string
	switch r.Status {
	case Pass:
		prefix = "PASS"
	case Fail:
		prefix = "FAIL"
	case Warn:
		prefix = "WARN"
	}
	_, _ = fmt.Fprintf(w, "%-4s  %-9s  %s\n", prefix, r.Group, r.Message)
	if r.Fix != "" {
		_, _ = fmt.Fprintf(w, "      fix: %s\n", r.Fix)
	}
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
