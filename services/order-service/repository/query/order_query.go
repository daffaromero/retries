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
	GetOrder(context.Context, *pb.GetOrderFilter, string) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest, pb.OrderService_GetOrdersServer) error
}

type OrderQueryImpl struct {
	Db *pgxpool.Pool
}

func (repo *OrderQueryImpl) CreateOrder(ctx context.Context, or *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	query := `INSERT INTO orders (id, customer_id, product_id, quantity) VALUES ($1, $2, $3, $4) RETURNING id`
	id := ""
	err := repo.Db.QueryRow(ctx, query, or.CustomerId, or.ProductId, or.Quantity).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.CreateOrderResponse{Id: or.Id, Status: true}, nil
}

func (repo *OrderQueryImpl) GetOrder(ctx context.Context, of *pb.GetOrderFilter, customerId string) (*pb.GetOrderResponse, error) {
	var CustomerId string
	query := `SELECT * from orders where customer_id=$1`
	err := repo.Db.QueryRow(ctx, query, of.CustomerId).Scan(&CustomerId)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("no records found")
	}
	if err != nil {
		return nil, fmt.Errorf("GetOrder: Bad input :: %e", err)
	} else {
		log.Println(of.CustomerId)
	}
	return &pb.GetOrderResponse{Orders: []*pb.Order{}}, nil
}

func (repo *OrderQueryImpl) GetAllOrders(ctx context.Context, req *pb.GetOrdersRequest, stream pb.OrderService_GetOrdersServer) error {
	query := `SELECT * from orders LIMIT $1 OFFSET $2`
	rows, err := repo.Db.Query(ctx, query, req.Count, req.Start)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		order := &pb.Order{}
		if err := rows.Scan(&order.Id, &order.CustomerId); err != nil {
			return err
		}
		response := &pb.GetOrderResponse{
			Orders: []*pb.Order{order},
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}
