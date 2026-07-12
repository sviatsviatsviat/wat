package main

import (
	"fmt"
	"os"

	_ "github.com/sviatsviatsviat/wat/agenthooks" // TODO: remove blank import when cmd/wat calls agenthooks APIs
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
			fmt.Fprint(os.Stdout, usage)
			return
		}
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
	fmt.Fprint(os.Stderr, usage)
	os.Exit(1)
}
