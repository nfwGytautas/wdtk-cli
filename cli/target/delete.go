package target

import (
	"os"
	"path"

	"github.com/nfwGytautas/mstk/cli/common"
	"github.com/nfwGytautas/mstk/cli/project"
	"github.com/urfave/cli"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Action for delete command
*/
func DeleteAction(ctx *cli.Context) {
	// Create bin directory
	cwd, err := os.Getwd()
	common.PanicOnError(err, "Failed to get current working directory")

	common.LogInfo("Deleting %s", path.Base(cwd))

	// Check if mstk_project.toml exists
	pc := project.ProjectConfig{}
	common.PanicOnError(pc.Read(), "Failed to read mstk_project.toml")

	// Teardown
	TeardownAction(ctx)

	// Remove namespace
	pc.Kubernetes.DeleteNamespace()

	// TODO: Delete docker images

	err = os.Chdir("../")
	common.PanicOnError(err, "Failed to change directory, the project directory needs to be deleted manually")

	err = os.RemoveAll(cwd)
	common.PanicOnError(err, "Failed to delete project directory, the project directory needs to be deleted manually")

	common.LogInfo("Done.")
}
