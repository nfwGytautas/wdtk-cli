package coordinator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

/*
Struct for keeping information about a shard
*/
type Shard struct {
	Name string `json:"Name"`
	URL  string `json:"URL"`
}

/*
Send a request to a coordinator requesting for a list of shards assigned to the service
*/
func GetShards(service string) []Shard {
	req, err := createCoordinatorRequest(http.MethodGet, "/shards/")
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

	var result []Shard
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.Println(err)
		return nil
	}

	return result
}
