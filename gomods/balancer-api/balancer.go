package balancer

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/nfwGytautas/mstk/gomods/common"
	"github.com/nfwGytautas/mstk/gomods/coordinator-api"
)

/*
Struct holding information for this balancer
*/
var BalancerInfo struct {
	mu sync.RWMutex

	Service     string              `toml:"Service"` // Service that the balancer is managing
	close       bool                // Flag for checking if the balancer should close
	ShardUpdate int                 `toml:"ShardUpdate"` // Number in ms how often to request for shard instance updates
	Shards      []coordinator.Shard // Available shards for the balancer
}

/*
Setup load balancer library
*/
func Setup() {
	log.Println("Setting up balancer lib")

	// Setup coordinator
	coordinator.Setup(os.Args[1])

	// Read balancer config
	err := common.StoreTOMLConfig(os.Args[2], &BalancerInfo)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Service: %s", BalancerInfo.Service)

	// Balancer open by default
	BalancerInfo.close = false

	go readAvailableShards()
}

/*
Routine continually reads available shards from a coordinator
*/
func readAvailableShards() {
	log.Printf("Querying new shards at speed of: %v", BalancerInfo.ShardUpdate)
	ticker := time.NewTicker(time.Duration(time.Duration(BalancerInfo.ShardUpdate) * time.Millisecond))

	for range ticker.C {
		if BalancerInfo.close {
			// Stop
			ticker.Stop()
			return
		}

		// Query coordinator for shards
		shards := coordinator.GetShards(BalancerInfo.Service)

		if shards != nil {
			log.Printf("Got shards: %v", shards)

			BalancerInfo.mu.Lock()
			BalancerInfo.Shards = shards
			BalancerInfo.mu.Unlock()
		}
	}
}
