package target

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

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
	EnsureMSTKRoot()

	// Remove kubectl
	services := GetMstkServicesList()

	var wg sync.WaitGroup
	wg.Add(len(services))

	for _, service := range services {
		go cleanupKubes(strings.ToLower(service), &wg)
	}

	wg.Wait()

	// Remove bin folder
	log.Println("Deleting bin")
	os.RemoveAll("./bin")
}

/*
Cleanup kubernetes
*/
func cleanupKubes(service string, wg *sync.WaitGroup) {
	defer TimeFn(fmt.Sprintf("Cleaning up %s", service))()
	defer wg.Done()

	log.Printf("Cleaning up %s", service)

	cleanupCmd := exec.Command("kubectl", "delete", "-f", fmt.Sprintf("kubes/%s/", service))
	log.Printf("Running %s", cleanupCmd.String())

	_, err := cleanupCmd.Output()
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
}
