package target

import (
	"log"
	"os"

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

	mstkDir, err := common.GetMSTKDir()
	common.PanicOnError(err, "Failed to get MSTK directory")

	// Remove kubectl from all projects
	// TODO: Project registry?
	// TODO: Delete mstk images

	// Remove bin folder
	log.Println("Deleting mstk directory")
	common.PanicOnError(os.RemoveAll(mstkDir), "Failed to remove mstk directory")
}

// PRIVATE FUNCTIONS
// ========================================================================
