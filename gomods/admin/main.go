package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common-api"
	"gorm.io/gorm"
)

// TODO: Security improvements

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

var serviceDbConn common.DatabaseConnection

const dcs_services = "mstk:mstk123@tcp(coordinator-db:3306)/coordinator_db?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	log.Println("Setting up MSTK system admin service")

	serviceDbConn = common.DatabaseConnection{}
	serviceDbConn.Initialize(common.DatabaseConnectionConfig{
		DCS: dcs_services,
		MigrateCallback: func(d *gorm.DB) {
			d.AutoMigrate(&Service{}, &Endpoint{})
		},
	})

	// Create gin engine
	r := gin.Default()

	gs := r.Group("/services", common.RequireDatabaseConnectionMiddleware(&serviceDbConn))
	gs.GET("/", func(c *gin.Context) {
		var services []Service
		result := serviceDbConn.DB.Preload("Endpoints").Preload("Shards").Find(&services)

		if result.Error != nil {
			log.Println(result.Error)
			c.String(http.StatusInternalServerError, "Failed to query")
			return
		}

		c.IndentedJSON(http.StatusOK, services)
	})
	gs.POST("/", func(c *gin.Context) {
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

		err := serviceDbConn.DB.Create(&s).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "registration success"})
	})

	// Run gin and block routine
	r.Run(":8080")
}
