package target

import (
	"errors"
	"log"
	"os"
	"path"

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
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Deleting %s", path.Base(cwd))

	// Check if mstk_project.toml exists
	_, err = os.OpenFile("mstk_project.toml", os.O_RDONLY, os.ModePerm)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("mstk_project.toml not found")
		panic(50)
	}

	pc := project.ProjectConfig{}
	pc.Read()

	// Teardown
	TeardownAction(ctx)

	// Remove namespace
	deleteNamespace(pc.Project)

	err = os.Chdir("../")
	if err != nil {
		log.Panic(err)
	}

	err = os.RemoveAll(cwd)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Done")
}
