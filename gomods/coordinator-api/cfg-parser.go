package coordinator

import (
	"errors"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

/*
Struct for holding a single coordinator config
*/
type coordinator struct {
	Host string
}

/*
Struct for holding coordinator config
*/
var config struct {
	Master coordinator
	Backup coordinator
}

/*
Read coordinator configuration
*/
func readConfig() error {
	log.Println("Loading coordinator config")

	// Check that the path to config file is passed
	if len(os.Args) < 2 {
		return errors.New("NO COORDINATOR CONFIGURATION FILE PASSED")
	}

	fileContents, err := os.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(fileContents), &config)
	return err
}
