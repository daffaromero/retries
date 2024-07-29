package subscribers

import (
	"encoding/json"
	"log"

	"github.com/daffaromero/retries/services/common/genproto/purchases"
	"github.com/daffaromero/retries/services/order-service/config"
	"github.com/nats-io/nats.go"
)

func ConsumeOrders(js nats.JetStreamContext) {
	_, err := js.Subscribe(config.SubjectNameOrderCreated, func(m *nats.Msg) {
		err := m.Ack()

		if err != nil {
			log.Println("could not ack", err)
			return
		}

		var orders purchases.Order
		err = json.Unmarshal(m.Data, &orders)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("consumer => subject: %s - order_id: %d - customer_id: %d, product_id: %d, quantity: %d\n", m.Subject, &orders.OrderId, orders.CustomerId, orders.ProductId, orders.Quantity)
	})

	if err != nil {
		log.Println("Failed when subscribing")
		return
	}
}
