package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nfwGytautas/webdev-tk/cli/scaffold"
	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/urfave/cli/v2"
)

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// Create cli.Command for scaffold
func ScaffoldCommand() *cli.Command {
	return &cli.Command{
		Flags:  []cli.Flag{},
		Name:   "scaffold",
		Usage:  "Check that the directory matches the wdtk.yml file configuration",
		Action: runScaffold,
	}
}

// PRIVATE FUNCTIONS
// ========================================================================

func runScaffold(ctx *cli.Context) error {
	// Read wdtk.yml
	cfg := types.WDTKConfig{}
	err := cfg.Read()
	if err != nil {
		return err
	}

	println("🧐  Checking...")

	stats := types.ServiceCheckStats{}

	actions := getScaffoldActions()
	for _, action := range actions {
		err := action(cfg, &stats)
		if err != nil {
			fmt.Println("❗️  " + err.Error())
			break
		}
	}

	print("--- ")
	fmt.Printf("%d created, ", stats.NumCreatedServices)
	fmt.Printf("%d modified, ", stats.NumModifiedServices)

	if len(stats.UnusedServices) > 0 {
		color.New(color.FgYellow).Printf("%d unused ", len(stats.UnusedServices))

	} else {
		fmt.Printf("0 unused ")
	}

	println("--- ")

	println("✅  Scaffolding done")

	return nil
}

func getScaffoldActions() []types.ScaffoldAction {
	return []types.ScaffoldAction{
		scaffold.CreateLocalServices,
		scaffold.PullGitServices,
		scaffold.GenerateConfigs,
		scaffold.WriteGoWork,
	}
}
