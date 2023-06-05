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

type generalConfig struct {
	GatewayIp string `json:"Gateway"`
}

type authConfig struct {
	generalConfig
	ConnectionString string `json:"ConnectionString"`
}

// PUBLIC FUNCTIONS
// ========================================================================

// Generates dynamic files and stores them in generated
func GenerateDynamics(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	err := createLocatorTable(cfg, stats)
	if err != nil {
		return err
	}

	err = createAuthConfig(cfg, stats)
	if err != nil {
		return err
	}

	err = createServiceConfig(cfg, stats)
	if err != nil {
		return err
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================
func createLocatorTable(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
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
				FullRequestIp: *serviceDeployment.IP + ":" + *serviceDeployment.Port,
			})
		}

		// Write
		file, err := json.MarshalIndent(ld, "", "    ")
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprintf("deploy/generated/LocatorTable_%s.json", deployment.Name), file, 0644)
		if err != nil {
			return err
		}

		// Create deployment script
		gatewayDeployment, err := cfg.GetFilledGatewayDeployment(deployment.Name)
		if err != nil {
			return err
		}

		err = createDeployment(gatewayDeployment)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDeployment(deployment types.DeploymentConfig) error {
	// TODO: Remote deploy

	gatewayDeploymentDirectory := strings.Replace(*deployment.DeployDir, "%serviceName", "Gateway", -1)
	authDeploymentDirectory := strings.Replace(*deployment.DeployDir, "%serviceName", "Authentication", -1)

	err := os.MkdirAll(gatewayDeploymentDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(authDeploymentDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	data := templates.WDTKDeployData{
		Deployment: deployment.Name,
		GatewayDir: gatewayDeploymentDirectory,
		AuthDir:    authDeploymentDirectory,
	}

	return file.WriteTemplate(fmt.Sprintf("deploy/unix/WDTK_%s.sh", deployment.Name), templates.LocalDeployWDTKTemplate, data)
}

func createAuthConfig(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	for _, authEntry := range cfg.Authentication.Entry {
		ac := authConfig{}

		gatewayConfig, err := cfg.GetFilledGatewayDeployment(authEntry.Name)
		if err != nil {
			return err
		}

		ac.GatewayIp = *gatewayConfig.IP + ":" + *gatewayConfig.Port
		ac.ConnectionString = authEntry.ConnectionString

		// Write
		file, err := json.MarshalIndent(ac, "", "    ")
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprintf("deploy/generated/AuthConfig_%s.json", authEntry.Name), file, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func createServiceConfig(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	for _, deployment := range cfg.Deployments {
		gc := generalConfig{}

		gatewayConfig, err := cfg.GetFilledGatewayDeployment(deployment.Name)
		if err != nil {
			return err
		}

		gc.GatewayIp = *gatewayConfig.IP + ":" + *gatewayConfig.Port

		// Write
		file, err := json.MarshalIndent(gc, "", "    ")
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprintf("deploy/generated/ServiceConfig_%s.json", deployment.Name), file, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
