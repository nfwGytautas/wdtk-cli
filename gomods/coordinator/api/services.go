package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common"
)

// TODO: Test only
var services = []common.Service{
	{
		Name: "Calculator",
		URL:  "http://localhost:7070/",
		Endpoints: []common.Endpoint{
			{
				Name: "Endpoint1",
			},
			{
				Name: "Endpoint2",
			},
		},
	},
	{
		Name: "TestService2",
		URL:  "http://localhost:7071/",
		Endpoints: []common.Endpoint{
			{
				Name: "Endpoint1",
			},
			{
				Name: "Endpoint2",
			},
		},
	},
}

func getServicesList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, services)
}

func getServiceEndpoints(c *gin.Context) {
	service := c.Query("service")

	for _, itService := range services {
		if itService.Name == service {
			c.IndentedJSON(http.StatusOK, itService.Endpoints)
			return
		}
	}

	// Service not found
	c.Status(http.StatusNotFound)
}

/*
Adds service locator specific gin routes
*/
func SetupServicesRoutes(r *gin.Engine) {
	locator := r.Group("/locator")

	locator.GET("/", getServicesList)
	locator.GET("/endpoints", getServiceEndpoints)
}
