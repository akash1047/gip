package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: gip <command> [args]")
		os.Exit(1)
	}

	var err error
	switch cmd := os.Args[1]; cmd {
	case "fetch":
		err = runFetch(os.Args[2:])
	case "doctor":
		err = runDoctor(os.Args[2:])
	case "mcp":
		err = runMCP(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "gip: unknown command %q\n", cmd)
		os.Exit(1)
	}

	if err != nil && !errors.Is(err, flag.ErrHelp) {
		fmt.Fprintf(os.Stderr, "gip: %v\n", err)
		os.Exit(1)
	}
}
