package service

import (
	"context"

	"github.com/daffaromero/retries/services/common/genproto/purchases"
)

var ordersDb = make([]*purchases.Order, 0)

type PurchaseService struct {
}

func NewPurchaseService() *PurchaseService {
	return &PurchaseService{}
}

func (s *PurchaseService) CreateOrder(ctx context.Context, o *purchases.Order) error {
	ordersDb = append(ordersDb, o)
	return nil
}

func (s *PurchaseService) GetOrders(ctx context.Context) []*purchases.Order {
	return ordersDb
}
