package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/balancer-api"
)

func main() {
	log.Println("Starting a OneToOne load balancer")

	// First setup balancer
	balancer.Start(configFunc, filterFunc)
}

func configFunc() balancer.BalancerInfo {
	return balancer.BalancerInfo{
		Service: "Calculator",
	}
}

func filterFunc(c *gin.Context) balancer.Shard {
	// Just forward to the first shard
	return balancer.Shard{
		URL: "shard-calculator:8080",
	}
}
