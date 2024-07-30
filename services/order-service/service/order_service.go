package service

import (
	"context"

	"github.com/daffaromero/retries/services/common/genproto/event"
	"github.com/daffaromero/retries/services/order-service/repository"
)

type OrderService interface {
	CreateOrder(context.Context, *event.EventRequest) (*event.EventResponse, error)
	GetOrder(context.Context, *event.GetEventFilter) (*event.GetEventResponse, error)
	GetAllOrders(context.Context, int, int)
}

type orderService struct {
	OrdRepo repository.OrderRepository
}
