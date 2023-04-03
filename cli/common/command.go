package common

import (
	"os/exec"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Execute a generic exec.Cmd command
*/
func ExecCmd(cmd *exec.Cmd) error {
	LogTrace("Running %s", cmd.String())
	output, err := cmd.Output()
	LogDebug("Output: %s", string(output))
	return err
}

// PRIVATE FUNCTIONS
// ========================================================================
