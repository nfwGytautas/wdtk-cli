package checks

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/templates"
	"github.com/nfwGytautas/webdev-tk/cli/types"
)

// PRIVATE TYPES
// ========================================================================
type locatorEntry struct {
	ServiceName   string `json:"Service"`
	FullRequestIp string `json:"IP"`
}

type locatorData struct {
	Mapping []locatorEntry `json:"Mapping"`
}

// PUBLIC FUNCTIONS
// ========================================================================

// Creates a locator table
func LocatorTableCreated(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("📍  Writing locator table")

	for _, deployment := range cfg.Deployments {
		ld := locatorData{}
		for _, service := range cfg.Services {
			serviceDeployment, err := cfg.GetFilledDeployment(service, deployment.Name)
			if err != nil {
				return err
			}

			ld.Mapping = append(ld.Mapping, locatorEntry{
				ServiceName:   service.Name,
				FullRequestIp: *serviceDeployment.IP,
			})
		}

		// Write
		file, err := json.MarshalIndent(ld, "", "    ")
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprintf("deploy/LD/%s.json", deployment.Name), file, 0644)
		if err != nil {
			return err
		}

		// Create deployment script
		gatewayDeployment, err := cfg.GetFilledGatewayDeployment(deployment.Name)
		if err != nil {
			return err
		}

		err = createGatewayDeployment(gatewayDeployment)
		if err != nil {
			return err
		}
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================
func createGatewayDeployment(deployment types.DeploymentConfig) error {
	// TODO: Remote deploy

	rootDeploymentDirectory := strings.Replace(*deployment.DeployDir, "%serviceName", "Gateway", -1)

	err := os.MkdirAll(rootDeploymentDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	data := templates.GatewayDeployData{
		Deployment: deployment.Name,
		OutDir:     rootDeploymentDirectory,
	}

	return file.WriteTemplate(fmt.Sprintf("deploy/unix/GATEWAY_%s.sh", deployment.Name), templates.LocalDeployGatewayTemplate, data)
}
