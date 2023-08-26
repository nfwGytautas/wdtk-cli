package types

import (
	"errors"
	"fmt"
	"strings"
)

type SourceConfig struct {
	Type     string  `yaml:"type"`
	Remote   *string `yaml:"remote,omitempty"`
	Language *string `yaml:"language,omitempty"`
}

type Options struct {
	IsGateway *bool `yaml:"gateway,omitempty"`
}

type ServiceDescriptionConfig struct {
	Name       string             `yaml:"name"`
	Source     SourceConfig       `yaml:"source"`
	Deployment []DeploymentConfig `yaml:"deployment"`
	Options    *Options           `yaml:"options,omitempty"`
}

func (service *ServiceDescriptionConfig) GitLocalDestination() (string, error) {
	if service.Source.Remote == nil {
		return "", errors.New("not a git service")
	}

	parts := strings.Split(*service.Source.Remote, "/")
	path := fmt.Sprintf("deploy/remotes/%s/%s", parts[2], parts[3])
	return path, nil
}
