// Wat is the wat hook command binary. Usage, supported commands, and flags are
// documented in [github.com/sviatsviatsviat/wat/internal/cli.RootHelpSummary].
package main

import (
	"os"

	"github.com/sviatsviatsviat/wat/internal/app"
)

func main() {
	os.Exit(app.Execute(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
