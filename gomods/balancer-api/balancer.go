package balancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Struct holding information for this balancer
*/
type BalancerInfo struct {
	Service  string             // Service that the balancer is managing
	filterFn LoadBalancerFilter // Filter to apply
}

type Shard struct {
	URL string
}

/*
Filter for balancing shards
*/
type LoadBalancerFilter func(*gin.Context) Shard

/*
Function for configuring a balancer
*/
type Configure func() BalancerInfo

/*
Configuration
*/
var balancerInfo BalancerInfo

/*
Start load balancer library
*/
func Start(config Configure, filter LoadBalancerFilter) {
	log.Println("Setting up balancer lib")

	// Get configuration
	balancerInfo = config()

	log.Printf("Service: %s", balancerInfo.Service)

	// Balancer open by default
	balancerInfo.filterFn = filter

	// Create gin engine
	r := gin.Default()

	// TODO: Rest of the endpoints
	r.GET("/*params", endpointHandler)
	r.POST("/*params", endpointHandler)

	// Run the balancer
	r.Run(":8080")
}

// ========================================================================
// PRIVATE
// ========================================================================

func endpointHandler(c *gin.Context) {
	shard := balancerInfo.filterFn(c)

	url, err := url.Parse("http://" + shard.URL)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("Forwarding '%s' -> '%s'", c.Request.URL.String(), url.String())

	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Director = func(req *http.Request) {
		req.Method = c.Request.Method
		req.Header = c.Request.Header
		req.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
