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
	Shards    []Shard
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
GORM shard struct
*/
type Shard struct {
	gorm.Model

	ServiceID uint

	Name  string
	URL   string
	State uint8
}

var dbConn common.DatabaseConnection

const dcs = "mstk:mstk123@tcp(coordinator-db:3306)/coordinator_db?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	log.Println("Setting up MSTK system admin service")

	dbConn = common.DatabaseConnection{}
	dbConn.Initialize(common.DatabaseConnectionConfig{
		DCS:             dcs,
		MigrateCallback: nil,
	})

	// Create gin engine
	r := gin.Default()

	gs := r.Group("/services", common.RequireDatabaseConnectionMiddleware(&dbConn))

	gs.GET("/", func(c *gin.Context) {
		var services []Service
		result := dbConn.DB.Preload("Endpoints").Preload("Shards").Find(&services)

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

		err := dbConn.DB.Create(&s).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "registration success"})
	})

	// Run gin and block routine
	r.Run(":8080")
}
