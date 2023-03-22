package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
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
func setupForwarding(r *gin.Engine, s coordinator.Service) error {
	log.Printf("Setting up %s", s.Name)

	log.Println(s.Name)
	for _, endpoint := range s.Endpoints {
		// Routing
		inUrl := fmt.Sprintf("/%s/%s", s.Name, endpoint.Name)
		outUrl := fmt.Sprintf("%s%s", s.URL, endpoint.Name)

		log.Printf("\t%s -> %s", inUrl, outUrl)

		// Create proxy handler
		handler, err := apiForward(outUrl)
		if err != nil {
			return err
		}

		// TODO: Other handles
		log.Println(inUrl)
		r.GET(inUrl, handler)
	}

	return nil
}
