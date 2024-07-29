package publishers

import (
	"github.com/daffaromero/retries/services/common/genproto/purchases"
	"github.com/nats-io/nats.go"
)

func publishPurchases(js nats.JetStreamContext) {

}

func getPurchases() ([]*purchases.Order, error) {
	allOrders, _ 
}
