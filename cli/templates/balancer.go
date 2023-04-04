package templates

// PUBLIC TYPES
// ========================================================================

/*
Data for balancer template
*/
type BalancerTemplateData struct {
	ServiceName string
}

/*
Template for balancer
*/
const BalancerTemplate = `
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/api/balancer-api"
)

type BalancerImplementation struct {
	// TODO: Add balancer data
}

func (bi *BalancerImplementation) GetServiceName() string {
	return "{{.ServiceName}}"
}

func (bi *BalancerImplementation) GetShard(ctx *gin.Context) (balancer.Shard, error) {
	// TODO: Add balancer filter
	return balancer.Shard{}, nil
}

func main() {
	balancer, err := balancer.CreateBalancer()
	if err != nil {
		log.Panicln("Failed to create balancer")
	}

	balancer.BalancerFn = &BalancerImplementation{}

	err = balancer.Run()
	if err != nil {
		log.Panic(err)
	}
}
`
