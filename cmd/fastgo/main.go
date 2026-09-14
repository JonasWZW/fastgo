// Command fastgo is the command-line entry point for the FastGo exercises.
package main

import (
	"os"

	"github.com/JonasWZW/fastgo/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
