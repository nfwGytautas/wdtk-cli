package coordinator

import (
	"fmt"
	"net/http"
)

// ========================================================================
// PRIVATE
// ========================================================================

/*
Make a request to the coordinator this will automatically resolve the correct URL
*/
func createCoordinatorRequest(method string, endpoint string) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("http://mstk-coordinator:8080%s", endpoint), nil)
}

/*
Execute a request to the coordinator
*/
func doRequest(r *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(r)
}
