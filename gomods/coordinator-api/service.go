package coordinator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

/*
Struct for holding information about a service
*/
type Service struct {
	Name      string     `json:"Name"`
	URL       string     `json:"URL"`
	Endpoints []Endpoint `json:"Endpoints"`
}

/*
Struct for holding information about a specific endpoint for a service
*/
type Endpoint struct {
	Name string `json:"Name"`
}

/*
Send a request to a coordinator requesting for a list of services
*/
func GetServices() []Service {
	req, err := createCoordinatorRequest(http.MethodGet, "/locator/")
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

	var result []Service
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.Println(err)
		return nil
	}

	return result
}
