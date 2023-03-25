package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
GORM service struct
*/
type Service struct {
	gorm.Model

	Name      string `gorm:"unique"`
	URL       string
	Endpoints []Endpoint
	Shards    []Shard
}

/*
GORM endpoint struct
*/
type Endpoint struct {
	gorm.Model

	ServiceID uint

	Name string
}

/*
GORM shard struct
*/
type Shard struct {
	gorm.Model

	ServiceID uint

	Name  string
	URL   string
	State uint8
}

func getServiceIdFromName(name string) uint {
	var s Service
	result := db.Where("name = ?", name).First(&s)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Println(result.Error)
		return 0
	}

	return s.ID
}

func getServicesList(c *gin.Context) {
	var services []Service
	result := db.Find(&services)

	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Failed to query")
		return
	}

	c.IndentedJSON(http.StatusOK, services)
}

func getServicesListExpanded(c *gin.Context) {
	var services []Service
	result := db.Preload("Endpoints").Preload("Shards").Find(&services)

	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Failed to query")
		return
	}

	c.IndentedJSON(http.StatusOK, services)
}

func getServiceEndpoints(c *gin.Context) {
	var s Service

	serviceName := c.Query("service")
	result := db.Where("name = ?", serviceName).Preload("Endpoints").First(&s)

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

func getServiceShards(c *gin.Context) {
	var s Service

	serviceName := c.Query("service")
	result := db.Where("name = ?", serviceName).Preload("Shards").First(&s)

	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Error retrieving")
		return
	}

	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.IndentedJSON(http.StatusOK, s.Shards)
}

func registerService(c *gin.Context) {
	// Request model
	input := struct {
		Name      string
		URL       string
		Endpoints []Endpoint
		Shards    []Shard
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := Service{}
	s.Name = input.Name
	s.URL = input.URL
	s.Endpoints = input.Endpoints
	s.Shards = input.Shards

	err := db.Create(&s).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func registerEndpoint(c *gin.Context) {
	// Request model
	input := struct {
		Name string
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
	e.ServiceID = getServiceIdFromName(serviceName)

	err := db.Create(&e).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func registerShard(c *gin.Context) {
	// Request model
	input := struct {
		Name  string
		URL   string
		State uint8
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

	s := Shard{}
	s.Name = input.Name
	s.URL = input.URL
	s.State = input.State
	s.ServiceID = getServiceIdFromName(serviceName)

	err := db.Create(&s).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

/*
Adds service locator specific gin routes
*/
func SetupServicesRoutes(r *gin.Engine) {
	locator := r.Group("/locator")

	// TODO: Authentication & Authorization
	locator.GET("/", getServicesList)
	locator.GET("/expanded", getServicesListExpanded)
	locator.GET("/endpoints", getServiceEndpoints)
	locator.GET("/shards", getServiceShards)

	locator.POST("/", registerService)
	locator.POST("/endpoints", registerEndpoint)
	locator.POST("/shards", registerShard)
}
