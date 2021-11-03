package main

import (
	"os"

	"github.com/Code-Hex/Neo-cowsay/internal/cli"
)

var version string

func main() {
	os.Exit((&cli.CLI{
		Version:  version,
		Thinking: true,
	}).Run())
}
