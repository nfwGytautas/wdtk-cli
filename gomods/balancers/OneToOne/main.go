package main

import (
	"log"
	"time"

	"github.com/nfwGytautas/mstk/gomods/balancer-api"
)

func main() {
	log.Println("Starting a OneToOne load balancer")

	// First setup balancer
	balancer.Setup()

	// Run the balancer
	time.Sleep(10 * time.Second)
}
