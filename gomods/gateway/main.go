package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/gateway/auth"
	"github.com/nfwGytautas/mstk/gomods/gateway/forward"
)

func main() {
	log.Println("Setting up API gateway")

	// Create gin engine
	r := gin.Default()

	// Setup authentication
	auth.Setup()

	// Configure gin
	auth.AddRoutes(r)

	// Configure forwarding routes
	forward.SetupRoutes(r)

	// Run gin and block routine
	r.Run(":8080")
}
