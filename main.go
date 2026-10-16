package main

import (
	"fmt"
	"os"

	"keyboard-tester/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "keyboard-tester:", err)
		os.Exit(1)
	}
}
