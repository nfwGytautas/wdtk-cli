package microservice

import (
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/lib/gdev/jwt"
)

// PUBLIC TYPES
// ========================================================================

/*
Interface for HTTP microservice implementation
*/
type HTTPMicroserviceImplementation interface {
	/*
		Function called when it is time to setup the microservice routes
	*/
	SetupRoutes(r *gin.Engine)
}

/*
Microservice using HTTP protocol
*/
type HTTPMicroservice struct {
	Implementation HTTPMicroserviceImplementation
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Create a HTTP microservice
*/
func CreateHTTPMicroservice() (HTTPMicroservice, error) {
	return HTTPMicroservice{}, nil
}

/*
Run the microservice

NOTE: Blocking goroutine
*/
func (ms *HTTPMicroservice) Run() error {
	r := gin.Default()
	r.Use(jwt.AuthenticationMiddleware())
	ms.Implementation.SetupRoutes(r)
	return r.Run(":8080")
}

// PRIVATE FUNCTIONS
// ========================================================================
