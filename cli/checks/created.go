package checks

import (
	"fmt"
	"os"

	"github.com/nfwGytautas/mstk/lib/gdev/array"
	"github.com/nfwGytautas/mstk/lib/gdev/file"
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

			err = os.Mkdir(path+"/service", os.ModePerm)
			if err != nil {
				return err
			}

			err = templates.WriteServiceTemplate(path + "/service/main.go")
			if err != nil {
				return err
			}

			if service.Source.Balancer != nil {
				err := os.Mkdir(path+"/balancer", os.ModePerm)
				if err != nil {
					return err
				}

				err = templates.WriteBalancerTemplate(path + "/balancer/main.go")
				if err != nil {
					return err
				}
			}

			stats.NumCreatedServices++
		} else {
			modified := false

			if !file.Exists(path + "/service") {
				err := os.Mkdir(path+"/service", os.ModePerm)
				if err != nil {
					return err
				}

				err = templates.WriteServiceTemplate(path + "/service/main.go")
				if err != nil {
					return err
				}

				modified = true
			}

			if service.Source.Balancer != nil {
				if !file.Exists(path + "/balancer") {
					err := os.Mkdir(path+"/balancer", os.ModePerm)
					if err != nil {
						return err
					}

					err = templates.WriteBalancerTemplate(path + "/balancer/main.go")
					if err != nil {
						return err
					}

					modified = true
				}
			}

			if modified {
				stats.NumModifiedServices++
			}
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
