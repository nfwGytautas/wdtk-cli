package templates

// PUBLIC TYPES
// ========================================================================

/*
Data of service template
*/
type ServiceTemplateData struct {
}

/*
Template for service main function
*/
const ServiceTemplate = `
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/backends/go/microservice-api"
)

type Microservice struct {
	// TODO: Add your microservice data
}

func (m *Microservice) SetupRoutes(r *gin.Engine) {
	// TODO: Setup routes
}

func main() {
	// TODO: Specify your microservice type
	microservice, err := microservice.CreateHTTPMicroservice()
	if err != nil {
		log.Panicln("Failed to create a service")
	}

	microservice.Implementation = &Microservice{}

	err = microservice.Run()
	if err != nil {
		log.Panic(err)
	}
}
`
