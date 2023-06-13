package checks

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nfwGytautas/webdev-tk/cli/types"
)

// PRIVATE TYPES
// ========================================================================
type locatorEntry struct {
	ServiceName   string `json:"service"`
	FullRequestIp string `json:"ip"`
}

// PUBLIC FUNCTIONS
// ========================================================================

// Generates dynamic files and stores them in generated
func GenerateDynamics(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	gatewayService, err := cfg.GetGatewayService()
	if err != nil {
		return err
	}

	for _, deployment := range cfg.Deployments {
		gatewayDeployment, err := cfg.GetFilledDeployment(gatewayService, deployment.Name)
		if err != nil {
			return err
		}

		// Services
		var locatorEntries []locatorEntry
		for _, service := range cfg.Services {
			if service.Name == gatewayService.Name {
				continue
			}

			serviceDeployment, err := cfg.GetFilledDeployment(service, deployment.Name)
			if err != nil {
				return err
			}

			locatorEntries = append(locatorEntries, locatorEntry{
				ServiceName:   service.Name,
				FullRequestIp: *serviceDeployment.IP + ":" + *serviceDeployment.Port,
			})

			configCopy := serviceDeployment.Config
			configCopy["runAddress"] = *serviceDeployment.IP + ":" + *serviceDeployment.Port
			configCopy["gatewayIp"] = *gatewayDeployment.IP + ":" + *gatewayDeployment.Port

			// Write
			file, err := json.MarshalIndent(configCopy, "", "    ")
			if err != nil {
				return err
			}

			err = os.WriteFile(fmt.Sprintf("deploy/generated/%s_ServiceConfig_%s.json", service.Name, deployment.Name), file, 0644)
			if err != nil {
				return err
			}
		}

		gatewayConfigCopy := gatewayDeployment.Config
		gatewayConfigCopy["runAddress"] = *gatewayDeployment.IP + ":" + *gatewayDeployment.Port
		gatewayConfigCopy["locatorTable"] = locatorEntries

		// Write gateway config
		file, err := json.MarshalIndent(gatewayConfigCopy, "", "    ")
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprintf("deploy/generated/Gateway_ServiceConfig_%s.json", deployment.Name), file, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
