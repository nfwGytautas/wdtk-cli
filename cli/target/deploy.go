package target

import (
	"log"

	"github.com/urfave/cli"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Action for deploy target
*/
func DeployAction(ctx *cli.Context) {
	defer TimeFn("Deploy")()

	log.Println("Deploying")

	// TODO: Implement
}
