package repository

import (
	"context"
	"fmt"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/order-service/repository/query"
	"github.com/daffaromero/retries/services/payment-service/processor"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v79"
)

type OrderRepository interface {
	CreateOrder(context.Context, *pb.Order) (*pb.Order, error)
	GetOrderDetails(context.Context, *pb.GetOrderFilter) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest) (*pb.GetOrderResponse, error)
	SendOrder(context.Context, *pb.SendOrderRequest) (*stripe.PaymentLink, error)
}

type orderRepository struct {
	db        Store
	ordQuery  query.OrderQuery
	processor processor.Stripe
}

func NewOrderRepository(db Store, ordQuery query.OrderQuery) OrderRepository {
	return &orderRepository{db: db, ordQuery: ordQuery}
}

func (o *orderRepository) CreateOrder(c context.Context, ord *pb.Order) (*pb.Order, error) {
	var res *pb.Order
	if err := o.db.WithTx(c, func(tx pgx.Tx) error {
		if _, err := o.ordQuery.GetOrderDetails(c, &pb.GetOrderFilter{CustomerId: ord.CustomerId}); err == nil {
			if err == pgx.ErrNoRows {
				re, err := o.ordQuery.CreateOrder(c, tx, ord)
				if err != nil {
					return err
				}
				res = re
			}
			if len(ord.ProductIds) == 0 {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderRepository) GetOrderDetails(c context.Context, fil *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	var res *pb.GetOrderResponse
	err := o.db.WithoutTx(c, func(pool *pgxpool.Pool) error {
		ord, err := o.ordQuery.GetOrderDetails(c, fil)
		if err != nil {
			return err
		}
		res = ord
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderRepository) GetAllOrders(c context.Context, req *pb.GetOrdersRequest) (*pb.GetOrderResponse, error) {
	var res *pb.GetOrderResponse
	err := o.db.WithoutTx(c, func(pool *pgxpool.Pool) error {
		ords, err := o.ordQuery.GetOrders(c, req)
		if err != nil {
			return err
		}
		res = ords
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderRepository) SendOrder(c context.Context, req *pb.SendOrderRequest) (*stripe.PaymentLink, error) {
	pl, err := o.processor.CreatePaymentLink(req)
	if err != nil {
		return nil, err
	}
	err = o.db.WithTx(c, func(tx pgx.Tx) error {
		err := o.ordQuery.SendOrder(c, tx, req, "pending", pl)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send order: %w", err)
	}

	return &stripe.PaymentLink{URL: pl}, nil
}
