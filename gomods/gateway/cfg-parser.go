package main

import (
	"errors"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

/*
Struct for holding gateway config
*/
var config struct {
	Port int
}

/*
Reads configuration file and loads the API gateway settings
*/
func readConfig() error {
	log.Println("Loading gateway config")

	// Check that the path to config file is passed
	if len(os.Args) < 3 {
		return errors.New("NO CONFIGURATION FILE PASSED")
	}

	fileContents, err := os.ReadFile(os.Args[2])
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(fileContents), &config)
	return err
}
