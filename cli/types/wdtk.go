package types

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// PUBLIC TYPES
// ========================================================================

type DeploymentConfig struct {
	Name        string  `yaml:"name"`
	IP          *string `yaml:"ip,omitempty"`
	BuildOnHost *bool   `yaml:"buildOnHost,omitempty"`
	DeployDir   *string `yaml:"dir,omitempty"`
	Port        *string `yaml:"port,omitempty"`
}

type ServiceConfig struct {
	Language string `yaml:"language"`
}

type BalancerConfig struct {
	Language string `yaml:"language"`
}

type ServiceDescriptionConfig struct {
	Name   string `yaml:"name"`
	Source struct {
		Service  *ServiceConfig  `yaml:"service"`
		Balancer *BalancerConfig `yaml:"balancer"`
	} `yaml:"source"`
	Deployment []DeploymentConfig `yaml:"deployment"`
}

// WDTK go representation of wdtk.yml file, generated with https://zhwt.github.io/yaml-to-go/
type WDTKConfig struct {
	Package     string             `yaml:"package"`
	Name        string             `yaml:"name"`
	Deployments []DeploymentConfig `yaml:"deployments"`
	APIGateway  struct {
		Deployment []DeploymentConfig `yaml:"deployment"`
	} `yaml:"apiGateway"`
	Services []ServiceDescriptionConfig `yaml:"services"`
}

// PRIVATE TYPES
// ========================================================================

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

	if serviceDeployment.BuildOnHost != nil {
		result.BuildOnHost = serviceDeployment.BuildOnHost
	}

	if serviceDeployment.Port != nil {
		result.Port = serviceDeployment.Port
	}

	return result, nil
}

// PRIVATE FUNCTIONS
// ========================================================================
