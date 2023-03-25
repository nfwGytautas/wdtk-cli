package coordinator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nfwGytautas/mstk/gomods/common"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Send a request to a coordinator requesting for a list of services
*/
func GetServices() []common.Service {
	req, err := createCoordinatorRequest(http.MethodGet, "/locator/expanded")
	if err != nil {
		log.Println(err)
		return nil
	}

	res, err := doRequest(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	var result []common.Service
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.Println(err)
		return nil
	}

	return result
}

/*
Get endpoints of a specific service
*/
func GetEndpoints(service string) []common.Endpoint {
	req, err := createCoordinatorRequest(http.MethodGet, "/locator/endpoints")
	if err != nil {
		log.Println(err)
		return nil
	}

	// Add service name
	q := req.URL.Query()
	q.Add("service", service)
	req.URL.RawQuery = q.Encode()

	res, err := doRequest(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	var result []common.Endpoint
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.Println(err)
		return nil
	}

	return result
}
