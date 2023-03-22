package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/coordinator/api"
)

func main() {
	log.Println("Setting up Coordinator")

	r := gin.Default()

	// Load config
	config, err := readConfig()
	if err != nil {
		log.Panic(err)
	}

	api.SetupServicesRoutes(r)

	// Run gin and block routine
	r.Run(config.Host)
}
