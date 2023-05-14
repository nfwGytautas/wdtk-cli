package checks

import (
	"fmt"
	"os"

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
func DeployScriptsExist(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("📦  Creating deployment scripts")

	for _, service := range cfg.Services {
		if !file.Exists(fmt.Sprintf("deploy/unix/%s.sh", service.Name)) {
			// Doesn't exist create
			stats.NumCreatedDeployScripts++

			err := createUNIXDeployScript(service)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================

func createUNIXDeployScript(service types.ServiceDescriptionConfig) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	balancerLang := ""
	serviceLang := ""

	if service.Source.Balancer != nil {
		balancerLang = service.Source.Balancer.Language
	}

	if service.Source.Service != nil {
		serviceLang = service.Source.Service.Language
	}

	data := templates.UNIXDeployData{
		ServiceName:  service.Name,
		RootDir:      currentDir,
		BalancerLang: balancerLang,
		ServiceLang:  serviceLang,
	}

	err = file.WriteTemplate(fmt.Sprintf("deploy/unix/%s.sh", service.Name), templates.UnixDeployTemplate, data)
	return err
}
