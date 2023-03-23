package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/balancer-api"
	"github.com/nfwGytautas/mstk/gomods/common"
)

func main() {
	log.Println("Starting a OneToOne load balancer")

	// First setup balancer
	balancer.Start(filterFunc)
}

func filterFunc(c *gin.Context, shards []common.Shard) common.Shard {
	// Just forward to the first shard
	return shards[0]
}
