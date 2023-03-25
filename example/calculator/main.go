package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/microservice-api"
)

func main() {
	log.Println("Starting calculator service")
	microservice.Start(setup)
}

func setup(r *gin.Engine) {
	r.GET("/Add", func(ctx *gin.Context) {
		log.Println("Echo test")
		ctx.Status(http.StatusNoContent)
	})
}
