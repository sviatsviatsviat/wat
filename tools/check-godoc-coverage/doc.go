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
