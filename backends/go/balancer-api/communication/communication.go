package communication

import (
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/backends/go/balancer-api/implementation"
)

// PUBLIC TYPES
// ========================================================================

/*
The inner communication implementation of a balancer
*/
type BalancerCommunication interface {
	/*
		Handle the request as per the communication type
	*/
	HandleRequest(c *gin.Context, shard implementation.Shard)
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// PRIVATE FUNCTIONS
// ========================================================================
