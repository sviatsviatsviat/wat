package main

import (
	"fmt"
	"os"

	_ "github.com/sviatsviatsviat/wat/agenthooks"
)

const usage = `wat — run and install Go hook scripts for coding agents

Usage:

	wat [command]

Commands will be added in later releases. Run with -h or help for this message.
`

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h", "--help", "help":
			fmt.Fprint(os.Stderr, usage)
			return
		}
	}
	fmt.Fprint(os.Stderr, usage)
	os.Exit(1)
}
