package balancer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

	Service     string            `toml:"Service"` // Service that the balancer is managing
	close       bool              // Flag for checking if the balancer should close
	ShardUpdate int               `toml:"ShardUpdate"` // Number in ms how often to request for shard instance updates
	Shards      []common.Shard    // Available shards for the balancer
	URL         string            `toml:"URL"` // URL for hosting the balancer
	Endpoints   []common.Endpoint // Endpoints for the service
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

	// Query endpoints
	log.Println("Querying endpoints")
	BalancerInfo.Endpoints = coordinator.GetEndpoints(BalancerInfo.Service)
	log.Printf("Got endpoints: %v", BalancerInfo.Endpoints)

	// Start monitoring shards
	go readAvailableShards()
	go monitorShardStatus()
}

/*
Routine continually reads available shards from a coordinator
*/
func readAvailableShards() {
	log.Printf("Querying new shards at speed of: %vms", BalancerInfo.ShardUpdate)
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

/*
Routine for monitoring status of shards
*/
func monitorShardStatus() {
	log.Printf("Querying shard state at speed of: %vms", BalancerInfo.ShardUpdate)
	ticker := time.NewTicker(time.Duration(time.Duration(BalancerInfo.ShardUpdate) * time.Millisecond))

	for range ticker.C {
		if BalancerInfo.close {
			// Stop
			ticker.Stop()
			return
		}

		// Query a shard for it's state
		BalancerInfo.mu.RLock()
		for _, shard := range BalancerInfo.Shards {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sstate", shard.URL), nil)
			if err != nil {
				log.Println(err)
				shard.Busy = true
				continue
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println(err)
				shard.Busy = true
				continue
			}

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Println(err)
				shard.Busy = true
				continue
			}

			var result bool
			err = json.Unmarshal(resBody, &result)
			if err != nil {
				log.Println(err)
				shard.Busy = true
				continue
			}

			shard.Busy = result
		}
		BalancerInfo.mu.RUnlock()
	}
}
