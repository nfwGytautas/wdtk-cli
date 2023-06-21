package types

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	SERVICE_TYPE_GIT    = "git"
	SERVICE_TYPE_BINARY = "bin"
	SERVICE_TYPE_LOCAL  = "src"
)

type DeploymentConfig struct {
	Name      string                 `yaml:"name"`
	IP        *string                `yaml:"ip,omitempty"`
	DeployDir *string                `yaml:"dir,omitempty"`
	Port      *string                `yaml:"port,omitempty"`
	Config    map[string]interface{} `yaml:"config,omitempty"`
	ApiKey    *string                `yaml:"apiKey,omitempty"`
}

// WDTK go representation of wdtk.yml file, generated with https://zhwt.github.io/yaml-to-go/
type WDTKConfig struct {
	Package     string                     `yaml:"package"`
	Name        string                     `yaml:"name"`
	Deployments []DeploymentConfig         `yaml:"deployments"`
	Services    []ServiceDescriptionConfig `yaml:"services"`
}

// Reads wdtk.yml file
func (wdtk *WDTKConfig) Read() error {
	println("🔍  Reading wdtk.yml")

	in, err := os.ReadFile("wdtk.yml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(in, wdtk)
}

// Returns all services whose type matches the one specified
func (wdtk *WDTKConfig) GetServicesOfType(t string) []ServiceDescriptionConfig {
	var result []ServiceDescriptionConfig
	for _, service := range wdtk.Services {
		if service.Source.Type == t {
			result = append(result, service)
		}
	}

	return result
}

// Returns the service that has gateway options set to true
func (wdtk *WDTKConfig) GetGatewayService() (ServiceDescriptionConfig, error) {
	for _, service := range wdtk.Services {
		if service.Options != nil && service.Options.IsGateway != nil && *service.Options.IsGateway {
			return service, nil
		}
	}

	// TODO: Verify that one exists
	return ServiceDescriptionConfig{}, errors.New("no gateway service provided")
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

	if serviceDeployment.ApiKey != nil {
		result.ApiKey = serviceDeployment.ApiKey
	}

	// TODO: Fill and override if defined
	if serviceDeployment.Config != nil {
		result.Config = serviceDeployment.Config
	} else {
		result.Config = make(map[string]interface{})
	}

	return result, nil
}
