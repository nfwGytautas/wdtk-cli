package gateway

import (
	"log"

	"github.com/gin-gonic/gin"
)

/*
Setup should be called first in almost all instances. This method will setup the api gateway functions for further use
*/
func Setup(r* gin.Engine) error {
	log.Println("Setting up API gateway")

	// Load config
	config, err := readConfig()
	if err != nil {
		return err
	}

	log.Printf("Config read: %v services", len(config.Services))

	// Services config loaded create gateway for request forwarding
	for _, service := range config.Services {
		setupForwarding(r, &service)
	}

	return nil
}
