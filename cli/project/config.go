package project

import (
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// ========================================================================
// PUBLIC
// ========================================================================

// TODO: Verify config
// TODO: Remove services

/*
Config that is contained inside mstk_project.toml file
*/
type ProjectConfig struct {
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
Read config file into the struct
*/
func (pc *ProjectConfig) Read() {
	b, err := os.ReadFile("mstk_project.toml")
	if err != nil {
		log.Printf("Failed to read mstk_project.toml")
		panic(50)
	}

	err = toml.Unmarshal(b, pc)
	if err != nil {
		log.Printf("Failed to unmarshal project config")
		panic(51)
	}
}

/*
Write to config file
*/
func (pc *ProjectConfig) Write() {
	pc.Version = cliVersion
	pc.GoVersion = "go 1.20"

	// Check if mstk_project.toml already exists, if not create it
	f, err := os.OpenFile("mstk_project.toml", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create mstk_project.toml")
		panic(50)
	}

	defer f.Close()

	b, err := toml.Marshal(pc)
	if err != nil {
		log.Printf("Failed to marshal project config")
		panic(51)
	}

	_, err = f.Write(b)
	if err != nil {
		log.Printf("Failed to write config to file")
		panic(52)
	}

	// Write go.work
	pc.writeGoWork()
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Write a go.work file from project config
*/
func (pc *ProjectConfig) writeGoWork() {
	// Check if mstk_project.toml already exists, if not create it
	f, err := os.OpenFile("go.work", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create go.work %v", err.Error())
		panic(50)
	}

	defer f.Close()

	f.WriteString("// Version\n")
	f.WriteString(pc.GoVersion)
	f.WriteString("\n\n")

	f.WriteString("// Workspaces\n")
	for _, service := range pc.Services {
		f.WriteString(fmt.Sprintf("use ./services/%s/service/\n", service.Name))
		f.WriteString(fmt.Sprintf("use ./services/%s/balancer/\n", service.Name))
		f.WriteString("\n")
	}
}

/*
Version of MSTK
*/
const cliVersion = "v0.0"
