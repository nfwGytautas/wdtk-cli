package checks

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/types"
)

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// Build scripts for WDTK modules
func PullWDTKServices(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("📕  Pulling WDTK services")

	abs, err := filepath.Abs("deploy/")
	if err != nil {
		return err
	}

	var outb, errb bytes.Buffer
	var cmd *exec.Cmd

	if !file.Exists(abs + "/wdtk-services/") {
		cmd = exec.Command("git", "clone", "https://github.com/nfwGytautas/wdtk-services.git")
		cmd.Dir = abs
	} else {
		cmd = exec.Command("git", "pull")
		cmd.Dir = abs + "/wdtk-services/"
	}

	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()

	if err != nil {
		fmt.Println(outb.String())
		fmt.Println(errb.String())
		fmt.Println(err.Error())
		return err
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================
