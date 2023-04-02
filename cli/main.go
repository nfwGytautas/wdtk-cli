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
			// MSTK
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

			// User
			{
				Name:   "template",
				Usage:  "Generate a template project",
				Action: target.TemplateAction,
			},
			{
				Name:   "delete",
				Usage:  "Delete project in the current dir",
				Action: target.DeleteAction,
			},
			{
				Name:  "service",
				Usage: "Create a new service for mstk project",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Add a new service to the current project",
						Action: target.CreateServiceAction,
					},
					{
						Name:   "remove",
						Usage:  "Remove a service from the current project",
						Action: target.RemoveServiceAction,
					},
				},
			},
			{
				Name:   "deploy",
				Usage:  "Deploy your mstk project to kubernetes",
				Action: target.DeployAction,
			},
			{
				Name:   "teardown",
				Usage:  "Teardown your mstk project from kubernetes",
				Action: target.TeardownAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
