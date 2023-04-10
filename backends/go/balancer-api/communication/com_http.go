package communication

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/backends/go/balancer-api/implementation"
)

// PUBLIC TYPES
// ========================================================================

/*
HTTP balancer communication implementation
*/
type HTTPBalancerCommunication struct {
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Handle the request.
HTTP to HTTP, forward the request via proxy
*/
func (bc *HTTPBalancerCommunication) HandleRequest(c *gin.Context, shard implementation.Shard) {
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

// PRIVATE FUNCTIONS
// ========================================================================
