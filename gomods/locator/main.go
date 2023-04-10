package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/services/locator/api"
	"github.com/nfwGytautas/mstk/gomods/services/locator/database"
)

func main() {
	log.Println("Setting up locator")

	// Setup database
	database.Setup()

	r := gin.Default()

	// Setup gin routes
	api.SetupServicesRoutes(r)

	// Run gin and block routine
	r.Run(":8080")
}
