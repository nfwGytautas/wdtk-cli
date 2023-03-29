package microservice

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Struct holding information for this microservice
*/
var Microservice struct {
	m sync.RWMutex

	Busy bool // State of the microservice true for busy, false otherwise
}

/*
Function for setting up a microservice
*/
type SetupMicroservice func(*gin.Engine)

/*
Setup microservice
*/
func Start(setupFn SetupMicroservice) {
	log.Println("Setting up microservice API")

	// Create gin engine, set it up, run
	r := gin.Default()

	addStateHandlers(r)
	setupFn(r)

	r.Run(":8080")
}

// ========================================================================
// PRIVATE
// ========================================================================

func addStateHandlers(r *gin.Engine) {
	g := r.Group("/state")

	g.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, Microservice.Busy)
	})
}
