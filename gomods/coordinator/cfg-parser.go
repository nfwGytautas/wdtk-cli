package main

import (
	"errors"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

/*
Struct for holding coordinator config
*/
type config struct {
	Name string
	Host string
}

/*
Read coordinator configuration
*/
func readConfig() (config, error) {
	log.Println("Loading coordinator config")
	cfg := config{}

	// Check that the path to config file is passed
	if len(os.Args) < 2 {
		return cfg, errors.New("NO CONFIGURATION FILE PASSED")
	}

	fileContents, err := os.ReadFile(os.Args[1])
	if err != nil {
		return cfg, err
	}

	_, err = toml.Decode(string(fileContents), &cfg)
	return cfg, err
}
