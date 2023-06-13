package types

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// PUBLIC TYPES
// ========================================================================

type AuthenticationEntry struct {
	DeploymentConfig `yaml:",inline"`
	ConnectionString string `yaml:"connectionString"`
}

type DeploymentConfig struct {
	Name      string  `yaml:"name"`
	IP        *string `yaml:"ip,omitempty"`
	DeployDir *string `yaml:"dir,omitempty"`
	Port      *string `yaml:"port,omitempty"`
}

type ServiceDescriptionConfig struct {
	Name       string             `yaml:"name"`
	Type       string             `yaml:"type"`
	Language   string             `yaml:"language"`
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
	Authentication struct {
		Entry []AuthenticationEntry `yaml:"deployment"`
	} `yaml:"authentication"`
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

	if serviceDeployment.Port != nil {
		result.Port = serviceDeployment.Port
	}

	return result, nil
}

// Get filled deployment data for api gateway
func (wdtk *WDTKConfig) GetFilledGatewayDeployment(deployment string) (DeploymentConfig, error) {
	var result DeploymentConfig
	var gatewayDeployment DeploymentConfig

	// Find the defined deployment
	for _, itDeployment := range wdtk.Deployments {
		if itDeployment.Name == deployment {
			result = itDeployment
		}
	}

	for _, itDeployment := range wdtk.APIGateway.Deployment {
		if itDeployment.Name == deployment {
			gatewayDeployment = itDeployment
		}
	}

	if result.Name == "" {
		return result, errors.New("deployment doesn't exist")
	}

	// Now override values
	if gatewayDeployment.IP != nil {
		result.IP = gatewayDeployment.IP
	}

	if gatewayDeployment.DeployDir != nil {
		result.DeployDir = gatewayDeployment.DeployDir
	}

	if gatewayDeployment.Port != nil {
		result.Port = gatewayDeployment.Port
	}

	return result, nil
}

// Get filled deployment data for authentication
func (wdtk *WDTKConfig) GetFilledAuthDeployment(deployment string) (DeploymentConfig, error) {
	var result DeploymentConfig
	var authDeployment AuthenticationEntry

	// Find the defined deployment
	for _, itDeployment := range wdtk.Deployments {
		if itDeployment.Name == deployment {
			result = itDeployment
		}
	}

	for _, itDeployment := range wdtk.Authentication.Entry {
		if itDeployment.Name == deployment {
			authDeployment = itDeployment
		}
	}

	if result.Name == "" {
		return result, errors.New("deployment doesn't exist")
	}

	// Now override values
	if authDeployment.IP != nil {
		result.IP = authDeployment.IP
	}

	if authDeployment.DeployDir != nil {
		result.DeployDir = authDeployment.DeployDir
	}

	if authDeployment.Port != nil {
		result.Port = authDeployment.Port
	}

	return result, nil
}

// PRIVATE FUNCTIONS
// ========================================================================
