package main

import (
	"log"
	"sync"

	"github.com/daffaromero/retries/services/purchases/service/publishers"
	"github.com/daffaromero/retries/services/purchases/service/subscribers"
	"github.com/daffaromero/retries/services/purchases/stream"
)

func main() {
	// httpServer := NewHTTPServer("localhost:9000")
	// go httpServer.Run()

	// gRPCServer := NewgRPCServer("localhost:8086")
	// gRPCServer.Run()

	js, err := stream.JetStreamInit()
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		publishers.PublishOrders(js)
	}()

	wg.Add(2)
	go func() {
		defer wg.Done()
		subscribers.ConsumeOrders(js)
	}()

	wg.Wait()
}
