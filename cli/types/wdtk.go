package types

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// PUBLIC TYPES
// ========================================================================

type DeploymentConfig struct {
	Name      string                 `yaml:"name"`
	IP        *string                `yaml:"ip,omitempty"`
	DeployDir *string                `yaml:"dir,omitempty"`
	Port      *string                `yaml:"port,omitempty"`
	Config    map[string]interface{} `yaml:"config,omitempty"`
}

type ServiceDescriptionConfig struct {
	Name       string             `yaml:"name"`
	Type       string             `yaml:"type"`
	Language   string             `yaml:"language"`
	Deployment []DeploymentConfig `yaml:"deployment"`
}

// WDTK go representation of wdtk.yml file, generated with https://zhwt.github.io/yaml-to-go/
type WDTKConfig struct {
	Package     string                     `yaml:"package"`
	Name        string                     `yaml:"name"`
	Deployments []DeploymentConfig         `yaml:"deployments"`
	Services    []ServiceDescriptionConfig `yaml:"services"`
}

// PRIVATE TYPES
// ========================================================================

const WDTK_SERVICE_IDENTIFIER = "wdtk"
const WDTK_GATEWAY_SERVICE_NAME = "Gateway"

// PUBLIC FUNCTIONS
// ========================================================================

// Reads wdtk.yml file
func (wdtk *WDTKConfig) Read() error {
	println("🔍  Reading wdtk.yml")

	in, err := os.ReadFile("wdtk.yml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(in, wdtk)
}

func (wdtk *WDTKConfig) GetGatewayService() (ServiceDescriptionConfig, error) {
	for _, service := range wdtk.Services {
		if service.Name == WDTK_GATEWAY_SERVICE_NAME {
			return service, nil
		}
	}

	return ServiceDescriptionConfig{}, errors.New("failed to get gateway service")
}

func (wdtk *WDTKConfig) GetUserServices() []ServiceDescriptionConfig {
	var result []ServiceDescriptionConfig
	for _, service := range wdtk.Services {
		if service.Type != WDTK_SERVICE_IDENTIFIER {
			result = append(result, service)
		}
	}
	return result
}

// Get a filled deployment for a specific service
func (wdtk *WDTKConfig) GetFilledDeployment(service ServiceDescriptionConfig, deployment string) (DeploymentConfig, error) {
	var result DeploymentConfig
	var serviceDeployment DeploymentConfig

	// Find the defined deployment
	for _, itDeployment := range wdtk.Deployments {
		if itDeployment.Name == deployment {
			result = itDeployment
		}
	}

	for _, itDeployment := range service.Deployment {
		if itDeployment.Name == deployment {
			serviceDeployment = itDeployment
		}
	}

	if result.Name == "" {
		return result, errors.New("deployment doesn't exist")
	}

	// Now override values
	if serviceDeployment.IP != nil {
		result.IP = serviceDeployment.IP
	}

	if serviceDeployment.DeployDir != nil {
		result.DeployDir = serviceDeployment.DeployDir
	}

	if serviceDeployment.Port != nil {
		result.Port = serviceDeployment.Port
	}

	// TODO: Fill and override if defined
	if serviceDeployment.Config != nil {
		result.Config = serviceDeployment.Config
	} else {
		result.Config = make(map[string]interface{})
	}

	return result, nil
}

// PRIVATE FUNCTIONS
// ========================================================================
