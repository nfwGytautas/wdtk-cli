package implementation

import "github.com/gin-gonic/gin"

// PUBLIC TYPES
// ========================================================================

/*
Interface for balancer functions
*/
type BalancerFunctions interface {
	/*
		Get the name of the balancer service
	*/
	GetServiceName() string

	/*
		Get the shard for appropriate request
	*/
	GetShard(ctx *gin.Context) (Shard, error)
}

/*
Struct representing a shard for the balancer
*/
type Shard struct {
	URL string
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// PRIVATE FUNCTIONS
// ========================================================================
