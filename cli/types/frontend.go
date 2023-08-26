package types

const (
	TOOLCHAIN_FLUTTER = "flutter"
)

type PlatformEntry struct {
	Type       string             `yaml:"type"`
	Toolchain  string             `yaml:"toolchain"`
	Deployment []DeploymentConfig `yaml:"deployment"`
}

type FrontendConfig struct {
	Platforms []PlatformEntry `yaml:"platforms"`
}

func (fc *FrontendConfig) GetFlutterPlatforms() []string {
	result := []string{}
	for _, entry := range fc.Platforms {
		if entry.Toolchain == TOOLCHAIN_FLUTTER {
			result = append(result, entry.Type)
		}
	}
	return result
}
