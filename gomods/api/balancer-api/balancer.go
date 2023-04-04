package balancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// PUBLIC TYPES
// ========================================================================

/*
A struct representing a balancer
*/
type Balancer struct {
	BalancerFn BalancerFunctions
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Creates a balancer and returns it or error if failed, make sure to set the balancers functions
*/
func CreateBalancer() (Balancer, error) {
	return Balancer{}, nil
}

/*
Run the balancer

NOTE: Blocking goroutine
*/
func (b *Balancer) Run() error {
	log.Printf("Starting balancer service for %s", b.BalancerFn.GetServiceName())

	r := gin.Default()

	// Routes
	r.GET("/*params", b.endpointHandler())
	r.POST("/*params", b.endpointHandler())

	return r.Run(":8080")
}

// PRIVATE FUNCTIONS
// ========================================================================

func (b *Balancer) endpointHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		shard, err := b.BalancerFn.GetShard(c)
		if err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}

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
}
