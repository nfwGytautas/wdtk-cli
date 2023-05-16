package balancer

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/gdev/jwt"
	"github.com/nfwGytautas/mstk/backends/go/balancer-api/communication"
	"github.com/nfwGytautas/mstk/backends/go/balancer-api/implementation"
)

// PUBLIC TYPES
// ========================================================================

/*
A struct representing a balancer
*/
type Balancer struct {
	Implementation implementation.BalancerFunctions
	Communication  communication.BalancerCommunication
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
	log.Printf("Starting balancer service for %s", b.Implementation.GetServiceName())

	r := gin.Default()

	// Input routes
	rootGroup := r.Group("/", jwt.AuthenticationMiddleware())
	rootGroup.Any("/*params", b.endpointHandler())

	return r.Run(":8080")
}

// PRIVATE FUNCTIONS
// ========================================================================

func (b *Balancer) endpointHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		shard, err := b.Implementation.GetShard(c)
		if err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		b.Communication.HandleRequest(c, shard)
	}
}
