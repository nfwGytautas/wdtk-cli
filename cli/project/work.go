package project

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Struct for managing a go work config file
*/
type GoWorkConfig struct {
	GoVersion     string
	UseDirectives []string
}

/*
Read go.work file
*/
func (gw *GoWorkConfig) Read() {
	file, err := os.Open("go.work")
	if err != nil {
		log.Printf("Failed to open go.work %v", err.Error())
		panic(50)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "go") {
			gw.GoVersion = scanner.Text()
		} else if strings.HasPrefix(scanner.Text(), "use") {
			gw.UseDirectives = append(gw.UseDirectives, scanner.Text())
		}
	}
}

/*
Write go work config to go.work
*/
func (gw *GoWorkConfig) Write() {
	gw.GoVersion = "go 1.20"

	// Check if mstk_project.toml already exists, if not create it
	f, err := os.OpenFile("go.work", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create go.work")
		panic(50)
	}

	defer f.Close()

	f.WriteString("// Version\n")
	f.WriteString(gw.GoVersion)
	f.WriteString("\n\n")

	f.WriteString("// Workspaces\n")
	for _, use := range gw.UseDirectives {
		f.WriteString("use ./" + use)
		f.WriteString("\n")
	}
}
