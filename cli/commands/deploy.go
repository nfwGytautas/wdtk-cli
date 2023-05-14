package commands

import (
	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/urfave/cli/v2"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

func DeployCommand() *cli.Command {
	return &cli.Command{
		Flags:     []cli.Flag{},
		Name:      "deploy",
		Usage:     "Deploy services",
		ArgsUsage: "[all|services...]",
		Action:    runDeploy,
	}
}

// PRIVATE FUNCTIONS
// ========================================================================

func runDeploy(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		println("❌  Deploy command expects either 'all', a service name or a list of services that you want to build")
		return nil
	}

	// Read wdtk.yml
	cfg := types.WDTKConfig{}
	err := cfg.Read()
	if err != nil {
		return err
	}

	println("✈️   Deploying...")

	return nil
}
