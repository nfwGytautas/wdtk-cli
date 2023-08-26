package scaffold

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/nfwGytautas/webdev-tk/cli/util"
)

type locatorEntry struct {
	ServiceName   string `json:"service"`
	FullRequestIp string `json:"ip"`
}

// PUBLIC FUNCTIONS
// ========================================================================
func GenerateConfigs(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("🏭  Generating service configurations")

	gatewayService, err := cfg.GetGatewayService()
	if err != nil {
		return err
	}

	for _, deployment := range cfg.Deployments {
		println(util.SPACING_1 + "- " + deployment.Name)
		err := generateDeployment(cfg, deployment, gatewayService)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateDeployment(cfg types.WDTKConfig, deployment types.DeploymentConfig, gateway types.ServiceDescriptionConfig) error {
	gatewayDeployment, err := cfg.GetFilledDeployment(gateway, deployment.Name)
	if err != nil {
		return err
	}

	var locatorEntries []locatorEntry
	for _, service := range cfg.Services {
		if service.Name == gateway.Name {
			continue
		}

		println(util.SPACING_2 + "- " + service.Name)

		serviceDeployment, err := cfg.GetFilledDeployment(service, deployment.Name)
		if err != nil {
			return err
		}

		locatorEntries = append(locatorEntries, locatorEntry{
			ServiceName:   service.Name,
			FullRequestIp: *serviceDeployment.IP + ":" + *serviceDeployment.Port,
		})

		err = generateServiceConfig(serviceDeployment, gatewayDeployment, service)
		if err != nil {
			return err
		}
	}

	println(util.SPACING_2 + "- " + gateway.Name)

	gatewayConfigCopy := gatewayDeployment.Config
	gatewayConfigCopy["runAddress"] = *gatewayDeployment.IP + ":" + *gatewayDeployment.Port
	gatewayConfigCopy["locatorTable"] = locatorEntries
	gatewayConfigCopy["apiKey"] = gatewayDeployment.ApiKey

	// Write gateway config
	file, err := json.MarshalIndent(gatewayConfigCopy, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("deploy/generated/Gateway_ServiceConfig_%s.json", deployment.Name), file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func generateServiceConfig(serviceDeployment, gatewayDeployment types.DeploymentConfig, service types.ServiceDescriptionConfig) error {
	configCopy := serviceDeployment.Config
	configCopy["runAddress"] = *serviceDeployment.IP + ":" + *serviceDeployment.Port
	configCopy["gatewayIp"] = *gatewayDeployment.IP + ":" + *gatewayDeployment.Port
	configCopy["apiKey"] = gatewayDeployment.ApiKey

	// Write
	file, err := json.MarshalIndent(configCopy, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("deploy/generated/%s_ServiceConfig_%s.json", service.Name, serviceDeployment.Name), file, 0644)
	if err != nil {
		return err
	}

	return nil
}
