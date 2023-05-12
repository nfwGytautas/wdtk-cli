package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nfwGytautas/webdev-tk/cli/checks"
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

	stats, err := serviceCheck(cfg)
	if err != nil {
		return err
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

	return nil
}

func serviceCheck(cfg types.WDTKConfig) (types.ServiceCheckStats, error) {
	stats := types.ServiceCheckStats{}
	var err error

	err = checks.AllServicesCreated(cfg, &stats)
	if err != nil {
		return stats, err
	}

	return stats, nil
}
