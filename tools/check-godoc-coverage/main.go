package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	root := flag.String("root", ".", "module root to scan")
	threshold := flag.Float64("threshold", 80, "minimum exported godoc coverage percent")
	listMissing := flag.Bool("list-missing", false, "print undocumented exported identifiers")
	flag.Parse()

	if math.IsNaN(*threshold) || *threshold < 0 || *threshold > 100 {
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
