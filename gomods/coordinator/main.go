package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common-api"
	"github.com/nfwGytautas/mstk/gomods/coordinator/api"
)

/*
Struct for holding coordinator config
*/
type config struct {
	Name string
	Host string
	DCS  string `toml:"DatabaseConnectionString"`
}

func main() {
	log.Println("Setting up Coordinator")

	r := gin.Default()

	// Load config
	cfg, err := common.ReadTOMLConfig[config](os.Args[1])
	if err != nil {
		log.Panic(err)
	}

	// Setup gin routes
	api.Setup(cfg.DCS)
	api.SetupServicesRoutes(r)

	// Run gin and block routine
	r.Run(cfg.Host)
}
