package common

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
