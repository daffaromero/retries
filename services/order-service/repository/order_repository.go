package repository

import (
	"context"

	"github.com/daffaromero/retries/services/common/genproto/event"
	"github.com/daffaromero/retries/services/order-service/repository/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(context.Context, *event.EventRequest) (*event.EventResponse, error)
	GetOrder(context.Context, *event.GetEventFilter) (*event.GetEventResponse, error)
	GetAllOrders(context.Context, int, int) ([]*event.GetEventResponse, error)
}

type orderRepository struct {
	db       Store
	ordQuery query.OrderQuery
}

func NewOrderRepository(db Store, ordQuery query.OrderQuery) OrderRepository {
	return &orderRepository{db: db, ordQuery: ordQuery}
}

func (o *orderRepository) CreateOrder(ctx context.Context, er *event.EventRequest) (*event.EventResponse, error) {
	var res *event.EventResponse
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		re, err := o.ordQuery.CreateOrder(ctx, er)
		if err != nil {
			return err
		}
		res = re
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderRepository) GetOrder(ctx context.Context, ef *event.GetEventFilter) (*event.GetEventResponse, error) {
	var res *event.GetEventResponse
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		re, err := o.ordQuery.GetOrder(ctx, ef)
		if err != nil {
			return err
		}
		res = re
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderRepository) GetAllOrders(ctx context.Context, count, start int) ([]*event.GetEventResponse, error) {
	var multires []*event.GetEventResponse
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		re, err := o.ordQuery.GetAllOrders(ctx, count, start)
		if err != nil {
			return err
		}
		multires = re
		return nil
	})
	if err != nil {
		return nil, err
	}
	return multires, nil
}
