package checks

import (
	"os"

	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/templates"
	"github.com/nfwGytautas/webdev-tk/cli/types"
)

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// Build scripts for WDTK modules
func WDTKBuild(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("📕  Creating WDTK modules build script")

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	data := templates.UNIXDeployData{
		RootDir: currentDir,
	}

	outFile := "deploy/unix/GATEWAY_BUILD_UNIX.sh"

	err = file.WriteTemplate(outFile, templates.UnixHeaderDeployTemplate, data)
	if err != nil {
		return err
	}

	err = file.AppendTemplate(outFile, templates.WDTKBuildTemplate, nil)
	if err != nil {
		return err
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================
