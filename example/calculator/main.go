package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/microservice-api"
)

func main() {
	log.Println("Starting calculator service")
	microservice.Start(setup)
}

func setup(r *gin.Engine) {

}
