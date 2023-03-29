package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/coordinator/api"
)

func main() {
	log.Println("Setting up Coordinator")

	r := gin.Default()

	// Setup gin routes
	api.Setup()
	api.SetupServicesRoutes(r)

	// Run gin and block routine
	r.Run(":8080")
}
