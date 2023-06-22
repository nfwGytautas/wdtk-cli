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
	Frontend    *FrontendConfig            `yaml:"frontend,omitempty"`
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
func (wdtk *WDTKConfig) GetFilledServiceDeployment(service ServiceDescriptionConfig, deployment string) (DeploymentConfig, error) {
	var serviceDeployment DeploymentConfig

	// Find the defined deployment
	result, err := wdtk.getDeploymentByName(deployment)
	if err != nil {
		return result, err
	}

	for _, itDeployment := range service.Deployment {
		if itDeployment.Name == deployment {
			serviceDeployment = itDeployment
		}
	}

	if serviceDeployment.Name == "" {
		return result, errors.New("deployment doesn't exist")
	}

	wdtk.fillDeployment(serviceDeployment, &result)
	return result, nil
}

func (wdtk *WDTKConfig) GetFilledFrontendDeployment(frontend PlatformEntry, deployment string) (DeploymentConfig, error) {
	var result DeploymentConfig
	var frontendDeployment DeploymentConfig

	// Find the defined deployment
	result, err := wdtk.getDeploymentByName(deployment)
	if err != nil {
		return result, err
	}

	for _, itDeployment := range frontend.Deployment {
		if itDeployment.Name == deployment {
			frontendDeployment = itDeployment
		}
	}

	if frontendDeployment.Name == "" {
		return result, errors.New("deployment doesn't exist")
	}

	wdtk.fillDeployment(frontendDeployment, &result)
	return result, nil
}

func (wdtk *WDTKConfig) getDeploymentByName(deployment string) (DeploymentConfig, error) {
	// Find the defined deployment
	for _, itDeployment := range wdtk.Deployments {
		if itDeployment.Name == deployment {
			return itDeployment, nil
		}
	}

	return DeploymentConfig{}, errors.New("failed to get deployment " + deployment)
}

func (wdtk *WDTKConfig) fillDeployment(source DeploymentConfig, target *DeploymentConfig) {
	// Now override values
	if source.IP != nil {
		target.IP = source.IP
	}

	if source.DeployDir != nil {
		target.DeployDir = source.DeployDir
	}

	if source.Port != nil {
		target.Port = source.Port
	}

	if source.ApiKey != nil {
		target.ApiKey = source.ApiKey
	}

	// TODO: Fill and override if defined
	if source.Config != nil {
		target.Config = source.Config
	} else {
		target.Config = make(map[string]interface{})
	}
}
