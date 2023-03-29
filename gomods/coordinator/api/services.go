package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common-api"
	"gorm.io/gorm"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
GORM service struct
*/
type Service struct {
	gorm.Model

	Name      string `gorm:"unique"`
	URL       string
	Endpoints []Endpoint
}

/*
GORM endpoint struct
*/
type Endpoint struct {
	gorm.Model

	ServiceID uint

	Name   string
	Method string
}

/*
Adds service locator specific gin routes
*/
func SetupServicesRoutes(r *gin.Engine) {
	locator := r.Group("/locator", common.RequireDatabaseConnectionMiddleware(&dbConn))

	// TODO: Authentication & Authorization
	locator.GET("/", getServicesList)
	locator.GET("/expanded", getServicesListExpanded)
	locator.GET("/service/expanded", getServiceExpanded)
	locator.GET("/endpoints", getServiceEndpoints)

	locator.POST("/", registerService)
	locator.POST("/endpoints", registerEndpoint)
}

// ========================================================================
// PRIVATE
// ========================================================================

func getServiceIdFromName(name string) uint {
	var s Service
	result := dbConn.DB.Where("name = ?", name).First(&s)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Println(result.Error)
		return 0
	}

	return s.ID
}

func getServicesList(c *gin.Context) {
	var services []Service
	result := dbConn.DB.Find(&services)

	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Failed to query")
		return
	}

	c.IndentedJSON(http.StatusOK, services)
}

func getServicesListExpanded(c *gin.Context) {
	var services []Service
	result := dbConn.DB.Preload("Endpoints").Find(&services)

	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Failed to query")
		return
	}

	c.IndentedJSON(http.StatusOK, services)
}

func getServiceExpanded(c *gin.Context) {
	serviceName := c.Query("service")
	if serviceName == "" {
		c.String(http.StatusBadRequest, "service not specified")
		return
	}

	var service Service
	result := dbConn.DB.Preload("Endpoints").Where("name = ?", serviceName).First(&service)

	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Failed to query")
		return
	}

	c.IndentedJSON(http.StatusOK, service)
}

func getServiceEndpoints(c *gin.Context) {
	var s Service

	serviceName := c.Query("service")
	result := dbConn.DB.Where("name = ?", serviceName).Preload("Endpoints").First(&s)

	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Error retrieving")
		return
	}

	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.IndentedJSON(http.StatusOK, s.Endpoints)
}

func registerService(c *gin.Context) {
	// Request model
	input := struct {
		Name      string
		URL       string
		Endpoints []Endpoint
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := Service{}
	s.Name = input.Name
	s.URL = input.URL
	s.Endpoints = input.Endpoints

	err := dbConn.DB.Create(&s).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func registerEndpoint(c *gin.Context) {
	// Request model
	input := struct {
		Name   string
		Method string
	}{}

	serviceName := c.Query("service")
	if serviceName == "" {
		c.String(http.StatusBadRequest, "No service provided")
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e := Endpoint{}
	e.Name = input.Name
	e.Method = input.Method
	e.ServiceID = getServiceIdFromName(serviceName)

	err := dbConn.DB.Create(&e).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}
