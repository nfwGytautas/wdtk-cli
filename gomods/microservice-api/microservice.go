package microservice

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
)

/*
Struct holding information for this microservice
*/
var Microservice struct {
	m sync.RWMutex

	URL  string
	Busy bool // State of the microservice true for busy, false otherwise
}

type SetupMicroservice func(*gin.Engine)

/*
Setup microservice
*/
func Start(setupFn SetupMicroservice) {
	log.Println("Setting up microservice API")

	// Setup coordinator API
	coordinator.Setup(os.Args[1])

	// Read microservice config
	err := common.StoreTOMLConfig(os.Args[2], &Microservice)
	if err != nil {
		log.Panic(err)
	}

	// Create gin engine, set it up, run
	r := gin.Default()

	addStateHandlers(r)
	setupFn(r)

	r.Run(Microservice.URL)
}

func addStateHandlers(r *gin.Engine) {
	g := r.Group("/state")

	g.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, Microservice.Busy)
	})
}
