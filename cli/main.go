package main

import (
	"log"
	"os"

	"github.com/nfwGytautas/mstk/cli/target"
	"github.com/urfave/cli"
)

// ========================================================================
// PUBLIC
// ========================================================================

func main() {
	app := &cli.App{
		Name:  "mstk",
		Usage: "CLI for MSTK",
		Commands: []cli.Command{
			{
				Flags:  target.SetupFlags,
				Name:   "setup",
				Usage:  "Setup mstk for usage (make sure kubectl is configured)",
				Action: target.SetupAction,
			},
			{
				Name:   "clean",
				Usage:  "Clean mstk (deletes bin/ folder as well as cleans up kubectl)",
				Action: target.CleanActionMstk,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
