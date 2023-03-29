package forward

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common-api"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Setup API gateway forwarding routes
*/
func SetupRoutes(r *gin.Engine) {
	gs := r.Group("/services")

	// TODO: Rest of the endpoints
	gs.GET("/:service/:endpoint/*params", handleRequest)
	gs.POST("/:service/:endpoint/*params", handleRequest)
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Returns true if the request can be forwarded, false otherwise
*/
func handleRequest(c *gin.Context) {
	// Get service name and endpoint
	serviceName := c.Param("service")
	endpointName := c.Param("endpoint")

	if serviceName == "" || endpointName == "" {
		c.String(http.StatusBadRequest, "service or endpoint not specified")
		return
	}

	// Check in with the service locator that the service name is valid
	service, err := coordinator.GetService(serviceName)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// We got the service now check if the endpoint is allowed
	var endpoint *common.Endpoint = nil
	for _, iEndpoint := range service.Endpoints {
		if endpointName == iEndpoint.Name {
			endpoint = &iEndpoint
		}
	}

	// Endpoint validity check
	if endpoint == nil {
		c.String(http.StatusBadRequest, "invalid endpoint")
		return
	}

	// TODO: Authentication and authorization

	// Endpoint and service valid proxy the request
	url, err := url.Parse(fmt.Sprintf("http://%s/%s", service.URL, endpoint.Name))
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "failed to create url for proxy")
		return
	}

	log.Printf("Forwarding '%s' -> '%s'", c.Request.URL.String(), url.String())

	// Create proxy and serve it
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Director = func(req *http.Request) {
		req.Method = c.Request.Method
		req.Header = c.Request.Header
		req.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = url.Path
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
