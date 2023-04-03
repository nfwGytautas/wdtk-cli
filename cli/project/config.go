package project

import (
	"fmt"
	"os"

	"github.com/nfwGytautas/mstk/cli/api"
	"github.com/pelletier/go-toml/v2"
)

// ========================================================================
// PUBLIC
// ========================================================================

// TODO: Verify config

/*
Data that is stored inside a .toml file
*/
type ProjectSaveData struct {
	Version         string         `toml:"MSTKVersion" comment:"The version of MSTK"`
	GoVersion       string         `toml:"GoVersion" comment:"Version of go to use"`
	Project         string         `toml:"Project" comment:"Project name"`
	PackageLocation string         `toml:"Package" comment:"Package location"`
	Services        []ServiceEntry `toml:"Services" comment:"Services of the project"`
}

/*
Entry for services inside mstk_project.toml
*/
type ServiceEntry struct {
	Name string `toml:"Name"`
}

/*
Config that is contained inside mstk_project.toml file
*/
type ProjectConfig struct {
	PSD        ProjectSaveData
	Kubernetes api.Kubernetes
	Builder    api.GoBuilder
	Docker     api.Docker
}

/*
Read config file into the struct
*/
func (pc *ProjectConfig) Read() error {
	b, err := os.ReadFile("mstk_project.toml")
	if err != nil {
		return err
	}

	err = toml.Unmarshal(b, &pc.PSD)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	pc.Kubernetes = api.CreateK8s(pc.PSD.Project)
	pc.Builder = api.CreateBuilder()
	pc.Docker = api.CreateDocker(cwd, pc.PSD.Project)
	return nil
}

/*
Write to config file
*/
func (pc *ProjectConfig) Write() error {
	pc.PSD.Version = cliVersion
	pc.PSD.GoVersion = "go 1.20"

	// Check if mstk_project.toml already exists, if not create it
	f, err := os.OpenFile("mstk_project.toml", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := toml.Marshal(pc.PSD)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	// Write go.work
	return pc.writeGoWork()
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Write a go.work file from project config
*/
func (pc *ProjectConfig) writeGoWork() error {
	// Check if mstk_project.toml already exists, if not create it
	f, err := os.OpenFile("go.work", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString("// Version\n")
	f.WriteString(pc.PSD.GoVersion)
	f.WriteString("\n\n")

	f.WriteString("// Workspaces\n")
	for _, service := range pc.PSD.Services {
		f.WriteString(fmt.Sprintf("use ./services/%s/service/\n", service.Name))
		f.WriteString(fmt.Sprintf("use ./services/%s/balancer/\n", service.Name))
		f.WriteString("\n")
	}

	return nil
}

/*
Version of MSTK
*/
const cliVersion = "v0.0"
