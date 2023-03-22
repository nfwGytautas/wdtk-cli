package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
)

func main() {
	log.Println("Setting up API gateway")

	// Setup coordinator API
	coordinator.Setup()

	// Create gin engine
	r := gin.Default()

	// Read config
	err := readConfig()
	if err != nil {
		log.Panic(err)
		return
	}

	// Forward services list
	for _, service := range coordinator.GetServices() {
		setupForwarding(r, service)
	}

	// Run gin and block routine
	r.Run(fmt.Sprintf("localhost:%v", config.Port))
}
