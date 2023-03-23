package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
)

/*
Struct for holding gateway config
*/
type config struct {
	Port int
}

func main() {
	log.Println("Setting up API gateway")

	// Setup coordinator API
	coordinator.Setup(os.Args[1])

	// Create gin engine
	r := gin.Default()

	// Read config
	cfg, err := common.ReadTOMLConfig[config](os.Args[2])
	if err != nil {
		log.Panic(err)
	}

	// Forward services list
	for _, service := range coordinator.GetServices() {
		setupForwarding(r, service)
	}

	// Run gin and block routine
	r.Run(fmt.Sprintf("localhost:%v", cfg.Port))
}
