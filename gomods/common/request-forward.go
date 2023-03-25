package common

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Function for generating a proxy forward for a url
*/
func ForwardRequestHandler(target string) (gin.HandlerFunc, error) {
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
