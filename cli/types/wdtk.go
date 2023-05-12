package types

import (
	"os"

	"gopkg.in/yaml.v3"
)

// PUBLIC TYPES
// ========================================================================

type DeploymentConfig struct {
	Name        string `yaml:"name"`
	IP          string `yaml:"ip"`
	BuildOnHost bool   `yaml:"buildOnHost"`
}

type ServiceConfig struct {
	Language string `yaml:"language"`
}

type BalancerConfig struct {
	Language string `yaml:"language"`
}

// WDTK go representation of wdtk.yml file, generated with https://zhwt.github.io/yaml-to-go/
type WDTKConfig struct {
	Name        string             `yaml:"name"`
	Deployments []DeploymentConfig `yaml:"deployments"`
	APIGateway  struct {
		Deployment []DeploymentConfig `yaml:"deployment"`
	} `yaml:"apiGateway"`
	Services []struct {
		Name   string `yaml:"name"`
		Source struct {
			Service  *ServiceConfig  `yaml:"service"`
			Balancer *BalancerConfig `yaml:"balancer"`
		} `yaml:"source"`
		Deployment []DeploymentConfig `yaml:"deployment"`
	} `yaml:"services"`
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

// PRIVATE FUNCTIONS
// ========================================================================
