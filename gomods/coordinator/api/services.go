package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
)

// TODO: Test only
var services = []coordinator.Service{
	{
		Name: "TestService1",
		URL:  "http://localhost:7070/TestServiceOne/",
		Endpoints: []coordinator.Endpoint{
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
		URL:  "http://localhost:7071/TestServiceTwo/",
		Endpoints: []coordinator.Endpoint{
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

/*
Adds service locator specific gin routes
*/
func SetupServicesRoutes(r *gin.Engine) {
	locator := r.Group("/locator")

	locator.GET("/", getServicesList)
}
