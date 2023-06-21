package scaffold

import (
	"fmt"
	"log"
	"os"

	"github.com/nfwGytautas/webdev-tk/cli/types"
)

// PUBLIC FUNCTIONS
// ========================================================================
func WriteGoWork(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("✏️   Writing go.work")

	err := resetGoWork()
	if err != nil {
		return err
	}

	for _, service := range cfg.GetServicesOfType(types.SERVICE_TYPE_LOCAL) {
		if *service.Source.Language == "go" {
			modPath := fmt.Sprintf("services/%s", service.Name)
			err := appendToGoWork(modPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Clear go.work and write the common boilerplate
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

// Append module to go.work
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
