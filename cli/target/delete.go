package target

import (
	"fmt"
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

	prompt, err := common.YNPrompt(fmt.Sprintf("Are you sure you want to delete the project %s", cwd), false)
	common.PanicOnError(err, "Failed prompt")
	if !prompt {
		return
	}

	common.LogInfo("Deleting %s", path.Base(cwd))

	// Check if mstk_project.toml exists
	pc := project.ProjectConfig{}
	common.PanicOnError(pc.Read(), "Failed to read mstk_project.toml")

	// Teardown
	TeardownAction(ctx)

	// Remove namespace
	pc.Kubernetes.DeleteNamespace()

	err = os.Chdir("../")
	common.PanicOnError(err, "Failed to change directory, the project directory needs to be deleted manually")

	err = os.RemoveAll(cwd)
	common.PanicOnError(err, "Failed to delete project directory, the project directory needs to be deleted manually")

	common.LogInfo("Done.")
}
