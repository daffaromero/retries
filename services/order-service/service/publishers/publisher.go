package publishers

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/daffaromero/retries/services/common/genproto/purchases"
	"github.com/daffaromero/retries/services/purchases/config"
	"github.com/nats-io/nats.go"
)

func PublishOrders(js nats.JetStreamContext) {
	orders, err := getPurchases()
	if err != nil {
		log.Println(err)
		return
	}

	for i := range orders {
		oneOrder := &orders[i]
		order := &oneOrder
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		orderString, err := json.Marshal(order)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = js.Publish(config.SubjectNameOrderCreated, orderString)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("publisher => message: %d\n", oneOrder.OrderId)
		}
	}
}

func getPurchases() ([]purchases.Order, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current working directory: %v", err)
		return nil, err
	}
	log.Printf("Current working directory: %s", wd)

	allOrders, err := os.ReadFile("./service/publishers/dummy.json")
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return nil, err
	}

	log.Printf("File content: %s", allOrders)

	var ordersObj []purchases.Order
	err = json.Unmarshal(allOrders, &ordersObj)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return nil, err
	}

	return ordersObj, nil
}
