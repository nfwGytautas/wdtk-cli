package scaffold

import (
	"fmt"
	"strings"

	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/nfwGytautas/webdev-tk/cli/util"
)

// PUBLIC FUNCTIONS
// ========================================================================
func CreateFlutterProject(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	if cfg.Frontend == nil {
		return nil
	}

	flutterPlatforms := cfg.Frontend.GetFlutterPlatforms()
	if len(flutterPlatforms) == 0 {
		return nil
	}

	fmt.Printf("📱  Creating flutter project for: %v\n", flutterPlatforms)

	platformString := strings.Join(flutterPlatforms, ",")

	err := util.ExecuteCommand(util.Command{
		Command:        "flutter",
		Args:           []string{"create", "_flutter", "--platforms=" + platformString},
		Directory:      "frontend/",
		PrintToConsole: false,
	})

	if err != nil {
		return err
	}

	return nil
}
