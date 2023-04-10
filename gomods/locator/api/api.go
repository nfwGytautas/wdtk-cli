package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/services/locator/database"
	"github.com/nfwGytautas/mstk/lib/gdev/jwt"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Adds service locator specific gin routes
*/
func SetupServicesRoutes(r *gin.Engine) {
	root := r.Group("/")

	root.Use(jwt.AuthenticationMiddleware(), jwt.AuthorizationMiddleware([]string{"_mstk"}))

	root.GET("/:service", getService)
	root.POST("/", registerService)
}

// PRIVATE FUNCTIONS
// ========================================================================

/*
Get a service
*/
func getService(c *gin.Context) {
	serviceName := c.Query("service")
	if serviceName == "" {
		c.String(http.StatusBadRequest, "service not specified")
		return
	}

	service, err := database.GetService(serviceName)

	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Failed to query")
		return
	}

	c.IndentedJSON(http.StatusOK, service)
}

/*
Register a service
*/
func registerService(c *gin.Context) {
	// Request model
	input := database.Service{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.CreateService(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}
