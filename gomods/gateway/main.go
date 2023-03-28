package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
	"github.com/nfwGytautas/mstk/gomods/gateway/auth"
)

func main() {
	log.Println("Setting up API gateway")

	// Setup coordinator API
	coordinator.Setup()

	// Create gin engine
	r := gin.Default()

	// Forward services list
	for _, service := range coordinator.GetServices() {
		setupForwarding(r, service)
	}

	// Setup authentication
	auth.Setup()

	// Configure gin
	auth.AddRoutes(r)

	// Run gin and block routine
	r.Run(":8080")
}
