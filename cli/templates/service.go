package templates

import "os"

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

	"github.com/nfwGytautas/wdtk-go-backend/microservice"
)

type Microservice struct {
	// TODO: Define the microservice
}

func main() {
	s, err := microservice.RegisterService(&Microservice{})
	if err != nil {
		log.Panicln("Failed to create a service")
	}

	// Specify microservice type
	s.CommunicationType = microservice.COMM_TYPE_HTTP

	err = s.Run()
	if err != nil {
		log.Panic(err)
	}
}
`

// Template for README.md in services directory
const ServicesReadME = `
# Services
Directory for all services
`

// PUBLIC FUNCTIONS
// ========================================================================
func WriteServiceTemplate(file string) error {
	return os.WriteFile(file, []byte(ServiceTemplate), os.ModePerm)
}
