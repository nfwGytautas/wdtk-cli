package coordinator

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/nfwGytautas/mstk/gomods/common"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Setup the coordinator API package
*/
func Setup(configFile string) {
	// Setup config
	err := common.StoreTOMLConfig(configFile, &config)
	if err != nil {
		log.Panic(err)
	}

	// Start monitoring the health of the coordinator
	go monitorCoordinatorHealth()

	// Wait for the monitor routine to start before returning
	time.Sleep(1 * time.Second)
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Struct for holding a single coordinator config
*/
type coordinator struct {
	Host string
}

/*
Struct for holding coordinator config
*/
var config struct {
	Master coordinator
	Backup coordinator
}

/*
Struct for keeping track of the coordinator state
*/
var cState struct {
	m         sync.RWMutex
	activeUrl string
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
	cState.activeUrl = config.Master.Host
}

/*
Make a request to the coordinator this will automatically resolve the correct URL
*/
func createCoordinatorRequest(method string, endpoint string) (*http.Request, error) {
	cState.m.RLock()
	defer cState.m.RUnlock()
	return http.NewRequest(method, fmt.Sprintf("%s%s", cState.activeUrl, endpoint), nil)
}

/*
Execute a request to the coordinator
*/
func doRequest(r *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(r)
}
