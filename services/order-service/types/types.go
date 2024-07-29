package types

import (
	"context"

	"github.com/daffaromero/retries/services/common/genproto/purchases"
)

type PurchaseService interface {
	CreateOrder(context.Context, *purchases.Order) error
	GetOrders(context.Context) []*purchases.Order
}
