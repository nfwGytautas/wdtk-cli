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
	"net/http"

	"github.com/nfwGytautas/wdtk-go-backend/microservice"
)

type Microservice struct {
	// TODO: Define the microservice context this can be accessed by every endpoint
}

func exampleEndpoint(e *microservice.EndpointExecutor) {
	e.Return(http.StatusOK, nil)
}

func main() {
	if microservice.RegisterService(microservice.ServiceDescription{
		ServiceContext: &Microservice{},
	}, []microservice.ServiceEndpoint{
		{
			Type:            microservice.ENDPOINT_TYPE_GET,
			Name:            "ExampleEndpoint/:id",
			Fn:              exampleEndpoint,
			EndpointContext: nil, // The struct that is passed here will be shared for every call to this endpoint
		},
	}) != nil {
		panic("Failed to run service")
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
