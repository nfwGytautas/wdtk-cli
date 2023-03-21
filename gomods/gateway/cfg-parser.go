package gateway

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

/*
 Struct for holding gateway config
*/
type config struct {
	Services map[string]service
}

/*
 Struct for holding a single service config
 */
type service struct {
	Name string
	Endpoints map[string]endpoint
}

/*
 Struct for holding information for an endpoint
 */
type endpoint struct {
	Name string
}

/*
 Reads configuration file "gateway.toml" and loads the API gateway settings
 */
func readConfig() (config, error) {
	log.Println("Loading gateway config")
	cfg := config{}

	fileContents, err := os.ReadFile("gateway.toml")
	if err != nil {
		return cfg, err
	}

	_, err = toml.Decode(string(fileContents), &cfg)
	return cfg, err
}
