package scaffold

import (
	"fmt"
	"strings"

	"github.com/nfwGytautas/gdev/array"
	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/nfwGytautas/webdev-tk/cli/util"
)

// PUBLIC FUNCTIONS
// ========================================================================
func PullGitServices(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("🪝   Acquiring git services")

	for _, pull := range getUniqueRemotes(cfg) {
		err := cloneOrPull(pull)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

// Parses all remotes that need to be pulled
func getUniqueRemotes(cfg types.WDTKConfig) []string {
	result := []string{}
	for _, service := range cfg.GetServicesOfType(types.SERVICE_TYPE_GIT) {
		parts := strings.Split(*service.Source.Remote, "/")
		path := strings.Join(parts[:3], "/")

		if !array.IsElementInArray(result, path) {
			result = append(result, path)
		}
	}

	return result
}

func cloneOrPull(remote string) error {
	repoName := strings.Split(remote, "/")[2]

	fmt.Print(util.SPACING_1 + "- ")

	if !file.Exists("deploy/remotes/" + repoName) {
		// Doesn't exist clone
		fmt.Printf("Cloning https://%s\n", remote)

		err := util.ExecuteCommand(util.Command{
			Command:   "git",
			Args:      []string{"clone", "https://" + remote},
			Directory: "deploy/remotes/",
		})

		if err != nil {
			return err
		}
	} else {
		// Exists pull
		fmt.Printf("Pulling https://%s\n", remote)

		err := util.ExecuteCommand(util.Command{
			Command:   "git",
			Args:      []string{"pull"},
			Directory: "deploy/remotes/" + repoName,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
