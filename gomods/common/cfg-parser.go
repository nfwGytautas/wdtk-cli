package common

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

/*
Read TOML from a specified config file and return a config file
*/
func ReadTOMLConfig[C any](configFile string) (C, error) {
	var cfg C
	err := StoreTOMLConfig(configFile, &cfg)
	return cfg, err
}

/*
Read TOML from a specified config file into an already constructed object
*/
func StoreTOMLConfig(configFile string, out interface{}) error {
	log.Printf("Loading config from %s into an object %T", configFile, out)

	fileContents, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(fileContents), out)
	return err
}
