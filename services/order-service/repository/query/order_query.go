package query

import (
	"context"
	"fmt"
	"log"

	pb "github.com/daffaromero/retries/services/common/genproto/orders"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderQuery interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	GetOrder(context.Context, *pb.GetOrderFilter) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest, pb.OrderService_GetOrdersServer) error
}

type OrderQueryImpl struct {
	Db *pgxpool.Pool
}

func NewOrderQueryImpl(db *pgxpool.Pool) *OrderQueryImpl {
	return &OrderQueryImpl{
		Db: db,
	}
}

func (o *OrderQueryImpl) CreateOrder(ctx context.Context, or *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	query := `INSERT INTO orders (id, customer_id, product_id, quantity) VALUES ($1, $2, $3, $4) RETURNING id`
	err := o.Db.QueryRow(ctx, query, or.Id, or.CustomerId, or.ProductId, or.Quantity).Scan(&or.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.CreateOrderResponse{Id: or.Id, Status: true}, nil
}

func (o *OrderQueryImpl) GetOrder(ctx context.Context, of *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	query := `SELECT * from orders WHERE customer_id=$1`
	rows, err := o.Db.Query(ctx, query, of.CustomerId)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("no records found")
	} else if err != nil {
		return nil, err
	}
	var orders []*pb.Order
	for rows.Next() {
		var order pb.Order
		err = rows.Scan(&order.Id, &order.CustomerId, &order.ProductId, &order.Quantity)
		if err != nil {
			return nil, err
		}
		newOrder := &pb.Order{
			Id:         order.Id,
			CustomerId: order.CustomerId,
			ProductId:  order.ProductId,
			Quantity:   order.Quantity,
		}
		orders = append(orders, newOrder)
	}
	return &pb.GetOrderResponse{Orders: orders}, nil
}

func (o *OrderQueryImpl) GetAllOrders(ctx context.Context, req *pb.GetOrdersRequest, stream pb.OrderService_GetOrdersServer) error {
	query := `SELECT id, customer_id, product_id, quantity FROM orders LIMIT $1 OFFSET $2`
	rows, err := o.Db.Query(ctx, query, req.Count, req.Start)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		order := &pb.Order{}

		if err := rows.Scan(&order.Id, &order.CustomerId, &order.ProductId, &order.Quantity); err != nil {
			return fmt.Errorf("scan error: %v", err)
		}
		response := &pb.GetOrderResponse{
			Orders: []*pb.Order{order},
		}
		err := stream.Send(response)
		if err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %v", err)
	}
	return nil
}
