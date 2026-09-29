package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: gip <command> [args]")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "gip: unknown command %q\n", os.Args[1])
	os.Exit(1)
}
