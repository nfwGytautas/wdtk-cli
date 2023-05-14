package main

import (
	"log"
	"os"

	"github.com/nfwGytautas/webdev-tk/cli/commands"
	"github.com/urfave/cli/v2"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// Version of the CLI
const Version = "0.0.0"

// PUBLIC FUNCTIONS
// ========================================================================

// PRIVATE FUNCTIONS
// ========================================================================

func main() {
	app := &cli.App{
		Name:                 "wdtk",
		EnableBashCompletion: true,
		Version:              Version,
		Commands: []*cli.Command{
			commands.InitCommand(),
			commands.ScaffoldCommand(),
			commands.DeployCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
