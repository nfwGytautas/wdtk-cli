package project

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Config that is contained inside mstk_project.toml file
*/
type ProjectConfig struct {
	Version         string         `toml:"MSTKVersion" comment:"The version of MSTK"`
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

	// Check if mstk_project.toml already exists, if not create it
	f, err := os.OpenFile("mstk_project.toml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Version of MSTK
*/
const cliVersion = "v0.0"
