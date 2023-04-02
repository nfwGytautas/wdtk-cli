package forward

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/api/common-api"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Setup API gateway forwarding routes
*/
func SetupRoutes(r *gin.Engine) {
	// Services
	gs := r.Group("/services")

	// Protect by authentication authorization needs to be done inside handleRequest
	gs.Use(common.AuthenticationMiddleware())

	// TODO: Rest of the endpoints
	gs.GET("/:service/:endpoint/*params", handleServicesRequest)
	gs.POST("/:service/:endpoint/*params", handleServicesRequest)

	// MSTK
	mstk := r.Group("/mstk")
	mstk.Use(common.AuthenticationMiddleware(), common.JwtAuthorizationMiddleware([]string{"_mstk"}))
	mstk.POST("/mstk/locator/:endpoint/*params", handleMstkLocatorRequest)
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Get a service with the specified name
*/
func getService(name string) (*common.Service, error) {
	req, err := http.NewRequest(http.MethodGet, "http://mstk-locator:8080/locator/", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("service", name)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result common.Service
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &result, nil
}

/*
Handles a /services/:service/* request
*/
func handleServicesRequest(c *gin.Context) {
	// Get service name and endpoint
	serviceName := c.Param("service")
	endpointName := c.Param("endpoint")

	if serviceName == "" || endpointName == "" {
		c.String(http.StatusBadRequest, "service or endpoint not specified")
		return
	}

	// Check in with the service locator that the service name is valid
	service, err := getService(serviceName)
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

	// Check if allowed role
	info, err := common.ParseToken(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Token")
		c.Abort()
	}

	if !common.IsElementInArray(endpoint.AllowedRoles, info.Role) {
		c.String(http.StatusUnauthorized, "Access denied")
		c.Abort()
		return
	}

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

/*
Handles a /mstk/locator/:endpoint/* request
*/
func handleMstkLocatorRequest(c *gin.Context) {
	endpointName := c.Param("endpoint")

	if endpointName == "" {
		c.String(http.StatusBadRequest, "service or endpoint not specified")
		return
	}

	// Endpoint and service valid proxy the request
	url, err := url.Parse(fmt.Sprintf("http://mstk-locator:8080/%s", endpointName))
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
