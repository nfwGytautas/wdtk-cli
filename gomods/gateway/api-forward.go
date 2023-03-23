package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common"
)

/*
Function sets up forwarding points for a given service
*/
func setupForwarding(r *gin.Engine, s common.Service) error {
	log.Printf("Setting up %s", s.Name)

	log.Println(s.Name)
	for _, endpoint := range s.Endpoints {
		// Routing
		inUrl := fmt.Sprintf("/%s/%s", s.Name, endpoint.Name)
		outUrl := fmt.Sprintf("%s%s", s.URL, endpoint.Name)

		log.Printf("\t%s -> %s", inUrl, outUrl)

		// Create proxy handler
		handler, err := common.ForwardRequestHandler(outUrl)
		if err != nil {
			return err
		}

		// TODO: Other handles
		log.Println(inUrl)
		r.GET(inUrl, handler)
	}

	return nil
}
