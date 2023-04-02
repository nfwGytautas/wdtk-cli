package target

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Clean action for mstk clean target
*/
func CleanActionMstk(ctx *cli.Context) {
	defer TimeFn("Cleaning")()

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	mstkDir := dirname + "/mstk/"

	// Remove kubectl
	cleanupCmd := exec.Command("kubectl", "delete", "-f", mstkDir)
	log.Printf("Running %s", cleanupCmd.String())

	_, err = cleanupCmd.Output()
	if err != nil {
		// Not found is not an actual error, just it doesn't exist which is fine since we are cleaning up anyway
		if !strings.Contains(
			string((err.(*exec.ExitError).Stderr)),
			"not found",
		) {
			log.Println(string((err.(*exec.ExitError).Stderr)))
			log.Panic(err)
		}
	}

	// Remove bin folder
	log.Println("Deleting mstk directory")
	os.Mkdir(mstkDir, os.ModePerm)
}
