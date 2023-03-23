package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
)

// TODO: Test only
var shards = map[string][]coordinator.Shard{
	"TestService1": {
		{
			Name: "Shard1",
			URL:  "http://localhost:6001/",
		},
		{
			Name: "Shard2",
			URL:  "http://localhost:6002/",
		},
	},
	"TestService2": {
		{
			Name: "Shard1",
			URL:  "http://localhost:6003/",
		},
		{
			Name: "Shard2",
			URL:  "http://localhost:6004/",
		},
		{
			Name: "Shard3",
			URL:  "http://localhost:6005/",
		},
	},
}

func getShardsList(c *gin.Context) {
	service := c.Query("service")
	log.Printf("Returning: %v", shards[service])
	c.IndentedJSON(http.StatusOK, shards[service])
}

/*
Adds service locator specific gin routes
*/
func SetupShardsRoutes(r *gin.Engine) {
	locator := r.Group("/shards")

	locator.GET("/", getShardsList)
}
