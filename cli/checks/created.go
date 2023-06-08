package checks

import (
	"fmt"
	"os"

	"github.com/nfwGytautas/gdev/array"
	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/templates"
	"github.com/nfwGytautas/webdev-tk/cli/types"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// Checks if all services are created that have been specified
func AllServicesCreated(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	serviceNames := []string{}

	for _, service := range cfg.Services {
		serviceNames = append(serviceNames, service.Name)
		path := fmt.Sprintf("services/%s", service.Name)

		// Check if root directory exists
		if !file.Exists(path) {
			// Create
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				return err
			}

			if service.Type == "service" {
				err := templates.WriteServiceTemplate(path + "/main.go")
				if err != nil {
					return err
				}
			} else if service.Type == "balancer" {
				err := templates.WriteBalancerTemplate(path + "/main.go")
				if err != nil {
					return err
				}
			}

			stats.NumCreatedServices++
		}
	}

	// Check if there are any excess
	dirs, err := file.GetDirectories("services/")
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !array.IsElementInArray(serviceNames, dir) {
			stats.UnusedServices = append(stats.UnusedServices, dir)
		}
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================
