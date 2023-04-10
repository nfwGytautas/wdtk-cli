package target

import (
	"os"

	"github.com/nfwGytautas/mstk/cli/api"
	"github.com/nfwGytautas/mstk/cli/common"
	"github.com/urfave/cli"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Clean action for mstk clean target
*/
func CleanActionMstk(ctx *cli.Context) {
	defer common.TimeCurrentFn()

	prompt, err := common.YNPrompt("Are you sure you want to delete MSTK?", false)
	common.PanicOnError(err, "Failed to prompt")
	if !prompt {
		return
	}

	mstkDir, err := common.GetMSTKDir()
	common.PanicOnError(err, "Failed to get MSTK directory")

	// Remove kubectl from all projects
	// TODO: Project registry?

	// Remove bin folder
	common.LogTrace("Deleting mstk directory")
	common.PanicOnError(os.RemoveAll(mstkDir), "Failed to remove mstk directory")

	prune, err := common.YNPrompt("Do you want to run 'docker container prune -a && docker image prune -a'?", false)
	common.PanicOnError(err, "Failed to prompt")
	if prune {
		common.PanicOnError(api.CleanDocker(), "Failed to clean docker")
	}

	common.LogInfo("Done.")
}

// PRIVATE FUNCTIONS
// ========================================================================
