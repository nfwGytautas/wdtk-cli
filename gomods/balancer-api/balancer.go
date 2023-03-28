package balancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common-api"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Struct holding information for this balancer
*/
type BalancerInfo struct {
	mu sync.RWMutex

	Service     string             // Service that the balancer is managing
	close       bool               // Flag for checking if the balancer should close
	ShardUpdate int                // Number in ms how often to request for shard instance updates
	Shards      []common.Shard     // Available shards for the balancer
	Endpoints   []common.Endpoint  // Endpoints for the service
	filterFn    LoadBalancerFilter // Filter to apply
}

/*
Filter for balancing shards
*/
type LoadBalancerFilter func(*gin.Context, []common.Shard) common.Shard

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

	// Setup coordinator
	coordinator.Setup()

	// Get configuration
	balancerInfo = config()

	log.Printf("Service: %s", balancerInfo.Service)

	// Balancer open by default
	balancerInfo.close = false
	balancerInfo.filterFn = filter

	// Create gin engine
	r := gin.Default()

	setupEndpoints(r)

	// Start monitoring shards
	go readAvailableShards()

	// Run the balancer
	r.Run("localhost:8080")
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Routine continually reads available shards from a coordinator
*/
func readAvailableShards() {
	log.Printf("Querying new shards at speed of: %vms", balancerInfo.ShardUpdate)
	ticker := time.NewTicker(time.Duration(time.Duration(balancerInfo.ShardUpdate) * time.Millisecond))

	for range ticker.C {
		if balancerInfo.close {
			// Stop
			ticker.Stop()
			return
		}

		// Query coordinator for shards
		shards := coordinator.GetShards(balancerInfo.Service)

		if shards != nil {
			log.Printf("Got shards: %v", shards)

			balancerInfo.mu.Lock()
			balancerInfo.Shards = shards
			balancerInfo.mu.Unlock()
		}
	}
}

func setupEndpoints(r *gin.Engine) {
	// Query endpoints
	log.Println("Querying endpoints")
	balancerInfo.Endpoints = coordinator.GetEndpoints(balancerInfo.Service)
	log.Printf("Got endpoints: %v", balancerInfo.Endpoints)

	// Create routes
	for _, endpoint := range balancerInfo.Endpoints {
		r.GET(endpoint.Name, endpointHandler)
	}
}

func endpointHandler(c *gin.Context) {
	if len(balancerInfo.Shards) == 0 {
		c.String(http.StatusPreconditionFailed, "No available shards")
		return
	}

	shard := balancerInfo.filterFn(c, balancerInfo.Shards)

	url, err := url.Parse("http://" + shard.URL)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
