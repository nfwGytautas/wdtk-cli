package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

/*
 Function for generating a proxy forward for a url
*/
func apiForward(target string) (gin.HandlerFunc, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(c *gin.Context) {
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = url.Host
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = url.Path
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}, nil
}

/*
 Function sets up forwarding points for a given service
 */
func setupForwarding(r* gin.Engine, s* service) error {
	log.Printf("Setting up %s", s.Name)

	for name, _ := range s.Endpoints {
		log.Printf("Routing endpoint %s", name)

		// Routing
		inUrl := fmt.Sprintf("/%s/%s", s.Name, name)
		outUrl := fmt.Sprintf("http://localhost:8080/rerouted/%s", name)

		// Create proxy handler
		handler, err := apiForward(outUrl)
		if err != nil {
			return err
		}

		// TODO: Other handles
		r.GET(inUrl, handler)
	}

	return nil
}
