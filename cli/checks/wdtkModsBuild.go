package checks

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

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

	outFile := "deploy/unix/WDTK_BUILD_UNIX.sh"

	err = file.WriteTemplate(outFile, templates.UnixHeaderDeployTemplate, data)
	if err != nil {
		return err
	}

	err = file.AppendTemplate(outFile, templates.WDTKBuildTemplate, nil)
	if err != nil {
		return err
	}

	abs, err := filepath.Abs("deploy/unix/")
	if err != nil {
		return err
	}

	println("🔨  Building WDTK services")

	// Run the deployment script
	var outb, errb bytes.Buffer

	cmd := exec.Command("bash", "./WDTK_BUILD_UNIX.sh")
	cmd.Dir = abs
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()

	// if err != nil {
	// 	file.Append(logFile, outb.String())
	// 	file.Append(logFile, errb.String())
	// 	file.Append(logFile, err.Error())

	// 	return err
	// } else {
	// 	file.Append(logFile, outb.String())
	// 	file.Append(logFile, errb.String())

	// 	return err
	// }

	return err
}

// PRIVATE FUNCTIONS
// ========================================================================
