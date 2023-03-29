package coordinator

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Setup the coordinator API package
*/
func Setup() {
	// Start monitoring the health of the coordinator
	go monitorCoordinatorHealth()

	// Wait for the monitor routine to start before returning
	time.Sleep(1 * time.Second)
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Struct for keeping track of the coordinator state
*/
var cState struct {
	m      sync.RWMutex
	online bool
}

/*
Function continually monitors the health of the current coordinator making sure it is online,
in the case that it is not the recipient URL will automatically be switched
*/
func monitorCoordinatorHealth() {
	// TODO: Implement coordinator health check
	cState.m.Lock()
	defer cState.m.Unlock()

	log.Println("Coordinator health monitor started")

	// Set to master by default
	cState.online = true
}

/*
Make a request to the coordinator this will automatically resolve the correct URL
*/
func createCoordinatorRequest(method string, endpoint string) (*http.Request, error) {
	cState.m.RLock()
	defer cState.m.RUnlock()

	if !cState.online {
		return nil, errors.New("coordinator offline")
	}

	return http.NewRequest(method, fmt.Sprintf("http://mstk-coordinator:8080/%s", endpoint), nil)
}

/*
Execute a request to the coordinator
*/
func doRequest(r *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(r)
}
