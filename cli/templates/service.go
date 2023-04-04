package templates

// PUBLIC TYPES
// ========================================================================

/*
Data of service template
*/
type ServiceTemplateData struct {
}

/*
Template for service main function
*/
const ServiceTemplate = `
package main

import "log"

func main() {
	log.Println("MSTK template service")
}

`
