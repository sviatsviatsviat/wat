// Command check-godoc-coverage reports and gates exported Go godoc coverage.
//
// Metric: fraction of exported funcs, methods, types, consts, and vars that
// have a doc comment (AST Doc on the declaration or group). Complements
// revive's exported rule (which still lint-fails undocumented exports,
// including interface methods). Struct/interface fields are not counted.
// Test files, testdata, and generated files (go.dev/s/generatedcode via
// ast.IsGenerated) are excluded.
//
// CodeRabbit's PR "docstring coverage" check is different: it scores the PR
// diff and may include unexported helpers. This tool is the CI source of truth
// for repository exported-godoc coverage (≥80%).
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	root := flag.String("root", ".", "module root to scan")
	threshold := flag.Float64("threshold", 80, "minimum exported godoc coverage percent")
	listMissing := flag.Bool("list-missing", false, "print undocumented exported identifiers")
	flag.Parse()

	if *threshold < 0 || *threshold > 100 {
		fmt.Fprintf(os.Stderr, "threshold must be between 0 and 100, got %v\n", *threshold)
		os.Exit(2)
	}

	report, err := measure(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "check-godoc-coverage: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("exported_godoc_coverage=%.2f%% documented=%d total=%d missing=%d threshold=%.2f%%\n",
		report.Percent, report.Documented, report.Total, len(report.Missing), *threshold)

	if *listMissing {
		for _, m := range report.Missing {
			fmt.Printf("MISSING: %s\n", m)
		}
	}

	if report.Total == 0 {
		fmt.Fprintln(os.Stderr, "check-godoc-coverage: no exported identifiers found")
		os.Exit(2)
	}
	if report.Percent+1e-9 < *threshold {
		fmt.Fprintf(os.Stderr, "exported godoc coverage %.2f%% is below threshold %.2f%%\n",
			report.Percent, *threshold)
		os.Exit(1)
	}
}
