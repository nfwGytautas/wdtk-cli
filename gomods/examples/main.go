package main

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/gateway"
)

func runReceiver(wg* sync.WaitGroup) {
	defer wg.Done()

	r := gin.Default()

	r.GET("/rerouted/test", func(ctx *gin.Context) {
		log.Println("Got rerouted request")
	})

	r.Run("localhost:8080")
}

func runGateway(wg* sync.WaitGroup) {
	defer wg.Done()

	r := gin.Default()

	err := gateway.Setup(r)
	if err != nil {
		panic(err)
	}

	r.Run("localhost:8081")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go runReceiver(&wg)
	go runGateway(&wg)

	wg.Wait()
}
