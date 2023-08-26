package types

type DeploymentConfig struct {
	Name      string                 `yaml:"name"`
	IP        *string                `yaml:"ip,omitempty"`
	DeployDir *string                `yaml:"dir,omitempty"`
	Port      *string                `yaml:"port,omitempty"`
	Config    map[string]interface{} `yaml:"config,omitempty"`
	ApiKey    *string                `yaml:"apiKey,omitempty"`
}

// Creates a deep copy of a deployment config
func (dc *DeploymentConfig) clone() DeploymentConfig {
	result := DeploymentConfig{}

	result.Name = dc.Name

	result.IP = pointerCopy(dc.IP)
	result.DeployDir = pointerCopy(dc.DeployDir)
	result.Port = pointerCopy(dc.Port)
	result.ApiKey = pointerCopy(dc.ApiKey)

	// Config map
	if dc.Config != nil {
		// Create the target map
		result.Config = make(map[string]interface{})

		// Copy from the original map to the target map
		for key, value := range dc.Config {
			result.Config[key] = value
		}
	}

	return result
}

func pointerCopy[T any](source *T) *T {
	if source == nil {
		return nil
	}

	value := new(T)
	*value = *source
	return value
}
