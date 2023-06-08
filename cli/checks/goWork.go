package checks

import (
	"fmt"
	"log"
	"os"

	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/templates"
	"github.com/nfwGytautas/webdev-tk/cli/types"
)

// PUBLIC FUNCTIONS
// ========================================================================

func GoWorkIsUpToDate(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("⚙️   Configuring packages, modules, etc.")
	err := resetGoWork()
	if err != nil {
		return err
	}

	for _, service := range cfg.Services {
		// Check go.mod
		if service.Language == "go" {
			modPath := fmt.Sprintf("services/%s", service.Name)
			goModPath := fmt.Sprintf("services/%s/go.mod", service.Name)

			if !file.Exists(goModPath) {
				err := writeGoMod(goModPath, cfg.Package, service.Name)
				if err != nil {
					return err
				}
			}

			err := appendToGoWork(modPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================

func resetGoWork() error {
	f, err := os.OpenFile("go.work", os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = f.WriteString("go 1.20\n")
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return f.Close()
}

func writeGoMod(path, pckg, name string) error {
	if !file.Exists(path) {
		data := templates.GoModFileData{
			Root:        pckg,
			ServiceName: name,
			GoVersion:   "1.20",
		}

		// Create go.mod
		err := file.WriteTemplate(path, templates.GoModTemplate, data)
		if err != nil {
			return nil
		}
	}

	return nil
}

func appendToGoWork(mod string) error {
	f, err := os.OpenFile("go.work",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if _, err := f.WriteString("use " + mod + "\n"); err != nil {
		log.Println(err)
	}

	return nil
}
