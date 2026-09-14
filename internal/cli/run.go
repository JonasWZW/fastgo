// Package cli implements the command-line boundary of the application.
package cli

import (
	"flag"
	"fmt"
	"io"

	"example.com/fastgo/internal/buildinfo"
)

// Run executes the CLI and returns a process exit code.
//
// Passing arguments and writers explicitly keeps the package independent from
// global process state and makes its behavior easy to test.
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("fastgo", flag.ContinueOnError)
	flags.SetOutput(stderr)

	showVersion := flags.Bool("version", false, "print build version")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "unexpected arguments: %v\n", flags.Args())
		return 2
	}

	if *showVersion {
		fmt.Fprintln(stdout, buildinfo.String())
		return 0
	}

	fmt.Fprintln(stdout, "fastgo: ready")
	return 0
}
