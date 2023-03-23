package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/balancer-api"
)

func main() {
	log.Println("Starting a OneToOne load balancer")

	// Create gin engine
	r := gin.Default()

	// First setup balancer
	balancer.Setup()

	// Run the balancer
	r.Run(balancer.BalancerInfo.URL)
}
