package api

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nfwGytautas/mstk/cli/common"
)

// PUBLIC TYPES
// ========================================================================

/*
Builder for go sources
*/
type GoBuilder struct {
	os   string
	arch string
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Create a GoBuilder
*/
func CreateBuilder() GoBuilder {
	// TODO: Find these automatically
	return GoBuilder{os: "linux", arch: "arm"}
}

/*
Build go sources
*/
func (gb *GoBuilder) Build(sourceDir, outputFile string) error {
	cmd := exec.Command("go", "build", "-o", outputFile, sourceDir)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", gb.os))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", gb.arch))

	return common.ExecCmd(cmd)
}

// PRIVATE FUNCTIONS
// ========================================================================
