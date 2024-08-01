package repository

import (
	"context"
	"log"

	pb "github.com/daffaromero/retries/services/common/genproto/orders"
	"github.com/daffaromero/retries/services/order-service/repository/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	GetOrder(context.Context, *pb.GetOrderFilter, string) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest, pb.OrderService_GetOrdersServer) error
}

type orderRepository struct {
	db       Store
	ordQuery query.OrderQuery
}

func NewOrderRepository(db Store, ordQuery query.OrderQuery) OrderRepository {
	return &orderRepository{db: db, ordQuery: ordQuery}
}

func (o *orderRepository) CreateOrder(ctx context.Context, ge *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Print(ge)
	var res *pb.CreateOrderResponse
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		re, err := o.ordQuery.CreateOrder(ctx, ge)
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

func (o *orderRepository) GetOrder(ctx context.Context, gf *pb.GetOrderFilter, customerId string) (*pb.GetOrderResponse, error) {
	var res *pb.GetOrderResponse
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		re, err := o.ordQuery.GetOrder(ctx, gf, customerId)
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

func (o *orderRepository) GetAllOrders(ctx context.Context, req *pb.GetOrdersRequest, sm pb.OrderService_GetOrdersServer) error {
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		return o.ordQuery.GetAllOrders(ctx, req, sm)
	})
	if err != nil {
		log.Print("the repo brokey")
		return err
	}
	return nil
}
