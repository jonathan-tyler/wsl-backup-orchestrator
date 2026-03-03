package main

import (
	"os"

	"github.com/jonathan-tyler/wsl-backup-orchestrator/internal/cli"
)

var (
	cliMain = cli.Main
	exitFunc = os.Exit
)

func main() {
	exitFunc(cliMain(os.Args[1:], os.Stdout, os.Stderr))
}
