package query

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderQuery interface {
	CreateOrder(c context.Context, tx pgx.Tx, order *pb.Order) (*pb.Order, error)
	GetOrderDetails(c context.Context, req *pb.GetOrderFilter) (*pb.GetOrderResponse, error)
	GetOrders(c context.Context, req *pb.GetOrdersRequest) (*pb.GetOrderResponse, error)
	UpdateOrder(c context.Context, tx pgx.Tx, order *pb.Order) (*pb.Order, error)
	SendOrder(c context.Context, tx pgx.Tx, req *pb.SendOrderRequest, status, paymentLink string) error
}

type OrderQueryImpl struct {
	db *pgxpool.Pool
}

func NewOrderQueryImpl(db *pgxpool.Pool) OrderQuery {
	return &OrderQueryImpl{db: db}
}

func (o *OrderQueryImpl) CreateOrder(c context.Context, tx pgx.Tx, req *pb.Order) (*pb.Order, error) {
	query := `INSERT INTO orders (id, customer_id, product_ids, products_details, settlement_status, total_payment, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := tx.Exec(c, query, req.Id, req.CustomerId, req.ProductIds, req.ProductsDetails, req.SettlementStatus, req.TotalPayment, req.CreatedAt, req.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.Order{
		Id:               req.Id,
		CustomerId:       req.CustomerId,
		ProductIds:       req.ProductIds,
		ProductsDetails:  req.ProductsDetails,
		SettlementStatus: req.SettlementStatus,
		TotalPayment:     req.TotalPayment,
		CreatedAt:        req.CreatedAt,
		UpdatedAt:        req.UpdatedAt,
	}, nil
}

func (o *OrderQueryImpl) GetOrderDetails(c context.Context, fil *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	query := `SELECT id, customer_id, product_ids, products_details, settlement_status, total_payment, created_at, updated_at from orders WHERE id=$1`
	rows, err := o.db.Query(c, query, fil.OrderId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no orders found for with ID %s", fil.OrderId)
		}
		log.Printf("Error querying orders: %v", err)
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}
	var orders []*pb.Order
	var order pb.Order
	err = rows.Scan(&order.Id, &order.CustomerId, &order.ProductIds, &order.ProductsDetails, &order.SettlementStatus, &order.TotalPayment, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	orders = append(orders, &order)
	return &pb.GetOrderResponse{Orders: orders}, nil
}
func (o *OrderQueryImpl) GetOrders(c context.Context, req *pb.GetOrdersRequest) (*pb.GetOrderResponse, error) {
	query := `SELECT id, customer_id, product_ids, products_details, settlement_status, total_payment, created_at, updated_at FROM orders LIMIT $1 OFFSET $2`
	rows, err := o.db.Query(c, query, req.Pagination.GetLimit(), req.Pagination.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var orders []*pb.Order
	for rows.Next() {
		var order pb.Order
		if err := rows.Scan(&order.Id, &order.CustomerId, &order.ProductIds, &order.ProductsDetails, &order.SettlementStatus, &order.TotalPayment, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return &pb.GetOrderResponse{Orders: orders}, nil
}

func (o *OrderQueryImpl) UpdateOrder(c context.Context, tx pgx.Tx, order *pb.Order) (*pb.Order, error) {
	query := `UPDATE orders SET customer_id = $1, product_ids = $2, products_details = $3, settlement_status = $4, total_payment = $5, updated_at = $6 WHERE id = $7 RETURNING id, customer_id, product_ids, products_details, settlement_status, total_payment, created_at, updated_at`
	var updatedOrder pb.Order
	err := tx.QueryRow(c, query, order.CustomerId, order.ProductIds, order.ProductsDetails, order.SettlementStatus, order.TotalPayment, order.UpdatedAt, order.Id).Scan(
		&updatedOrder.Id, &updatedOrder.CustomerId, &updatedOrder.ProductIds, &updatedOrder.ProductsDetails, &updatedOrder.SettlementStatus, &updatedOrder.TotalPayment, &updatedOrder.CreatedAt, &updatedOrder.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return &updatedOrder, nil
}

func (o *OrderQueryImpl) SendOrder(c context.Context, tx pgx.Tx, req *pb.SendOrderRequest, status, paymentLink string) error {
	query := `UPDATE orders SET settlement_status = $1, payment_link = $2, updated_at = $3 WHERE id = $3`
	_, err := tx.Exec(c, query, status, time.Now(), req.OrderId)
	if err != nil {
		return fmt.Errorf("failed to send order: %v", err)
	}
	return nil
}
