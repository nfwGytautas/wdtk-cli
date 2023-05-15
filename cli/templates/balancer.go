package templates

import "os"

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
	"github.com/nfwGytautas/webdev-tk/backends/go/balancer-api"
	"github.com/nfwGytautas/webdev-tk/backends/go/balancer-api/communication"
	"github.com/nfwGytautas/webdev-tk/backends/go/balancer-api/implementation"
)

type BalancerImplementation struct {
	// TODO: Add balancer data
}

func (bi *BalancerImplementation) GetServiceName() string {
	return "ExampleService"
}

func (bi *BalancerImplementation) GetShard(ctx *gin.Context) (implementation.Shard, error) {
	// TODO: Add balancer filter
	return implementation.Shard{}, nil
}

func main() {
	balancer, err := balancer.CreateBalancer()
	if err != nil {
		log.Panicln("Failed to create balancer")
	}

	balancer.Implementation = &BalancerImplementation{}

	// TODO: Specify your communication type
	balancer.Communication = &communication.HTTPBalancerCommunication{}

	err = balancer.Run()
	if err != nil {
		log.Panic(err)
	}
}

`

// PUBLIC FUNCTIONS
// ========================================================================
func WriteBalancerTemplate(file string) error {
	return os.WriteFile(file, []byte(BalancerTemplate), os.ModePerm)
}