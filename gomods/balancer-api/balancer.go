package balancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
var BalancerInfo struct {
	mu sync.RWMutex

	Service     string             `toml:"Service"` // Service that the balancer is managing
	close       bool               // Flag for checking if the balancer should close
	ShardUpdate int                `toml:"ShardUpdate"` // Number in ms how often to request for shard instance updates
	Shards      []common.Shard     // Available shards for the balancer
	URL         string             `toml:"URL"` // URL for hosting the balancer
	Endpoints   []common.Endpoint  // Endpoints for the service
	filterFn    LoadBalancerFilter // Filter to apply
}

/*
Filter for balancing shards
*/
type LoadBalancerFilter func(*gin.Context, []common.Shard) common.Shard

/*
Start load balancer library
*/
func Start(filter LoadBalancerFilter) {
	log.Println("Setting up balancer lib")

	// Setup coordinator
	coordinator.Setup(os.Args[1])

	// Read balancer config
	err := common.StoreTOMLConfig(os.Args[2], &BalancerInfo)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Service: %s", BalancerInfo.Service)

	// Balancer open by default
	BalancerInfo.close = false
	BalancerInfo.filterFn = filter

	// Create gin engine
	r := gin.Default()

	setupEndpoints(r)

	// Start monitoring shards
	go readAvailableShards()

	// Run the balancer
	r.Run(BalancerInfo.URL)
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Routine continually reads available shards from a coordinator
*/
func readAvailableShards() {
	log.Printf("Querying new shards at speed of: %vms", BalancerInfo.ShardUpdate)
	ticker := time.NewTicker(time.Duration(time.Duration(BalancerInfo.ShardUpdate) * time.Millisecond))

	for range ticker.C {
		if BalancerInfo.close {
			// Stop
			ticker.Stop()
			return
		}

		// Query coordinator for shards
		shards := coordinator.GetShards(BalancerInfo.Service)

		if shards != nil {
			log.Printf("Got shards: %v", shards)

			BalancerInfo.mu.Lock()
			BalancerInfo.Shards = shards
			BalancerInfo.mu.Unlock()
		}
	}
}

func setupEndpoints(r *gin.Engine) {
	// Query endpoints
	log.Println("Querying endpoints")
	BalancerInfo.Endpoints = coordinator.GetEndpoints(BalancerInfo.Service)
	log.Printf("Got endpoints: %v", BalancerInfo.Endpoints)

	// Create routes
	for _, endpoint := range BalancerInfo.Endpoints {
		r.GET(endpoint.Name, endpointHandler)
	}
}

func endpointHandler(c *gin.Context) {
	if len(BalancerInfo.Shards) == 0 {
		c.String(http.StatusPreconditionFailed, "No available shards")
		return
	}

	shard := BalancerInfo.filterFn(c, BalancerInfo.Shards)

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
