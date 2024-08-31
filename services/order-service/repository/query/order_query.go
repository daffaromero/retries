package query

import (
	"context"
	"fmt"
	"log"
<<<<<<< HEAD
	"time"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
=======

	pb "github.com/daffaromero/retries/services/common/genproto/orders"
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderQuery interface {
<<<<<<< HEAD
	CreateOrder(context.Context, pgx.Tx, *pb.Order) (*pb.Order, error)
	GetOrder(context.Context, *pb.GetOrderFilter) (*pb.GetOrderResponse, error)
	GetOrders(context.Context, *pb.GetOrdersRequest, pb.OrderService_GetOrdersServer) error
	UpdateOrder(context.Context, pgx.Tx, *pb.Order) (*pb.Order, error)
	SendOrder(context.Context, pgx.Tx, *pb.SendOrderRequest) error
}

type OrderQueryImpl struct {
	db *pgxpool.Pool
=======
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	GetOrder(context.Context, *pb.GetOrderFilter) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest, pb.OrderService_GetOrdersServer) error
}

type OrderQueryImpl struct {
	Db *pgxpool.Pool
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
}

func NewOrderQueryImpl(db *pgxpool.Pool) *OrderQueryImpl {
	return &OrderQueryImpl{
<<<<<<< HEAD
		db: db,
	}
}

func (o *OrderQueryImpl) CreateOrder(c context.Context, tx pgx.Tx, req *pb.Order) (*pb.Order, error) {
	query := `INSERT INTO orders (id, user_id, product_ids, products_details, total_payment, settlement_status, is_private, is_private_approved, private_seller_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	err := tx.QueryRow(c, query, req.Id, req.UserId, req.ProductIds, req.ProductsDetails, req.TotalPayment, req.SettlementStatus, req.SettlementStatus, req.IsPrivate, req.IsPrivateApproved, req.PrivateSellerId, req.CreatedAt, req.UpdatedAt).Scan(&req.Id, &req.UserId, &req.ProductIds, &req.ProductsDetails, &req.TotalPayment, &req.SettlementStatus, &req.SettlementStatus, &req.IsPrivate, &req.IsPrivateApproved, &req.PrivateSellerId, &req.CreatedAt, &req.UpdatedAt)
=======
		Db: db,
	}
}

func (o *OrderQueryImpl) CreateOrder(ctx context.Context, or *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	query := `INSERT INTO orders (id, customer_id, product_id, quantity) VALUES ($1, $2, $3, $4) RETURNING id`
	err := o.Db.QueryRow(ctx, query, or.Id, or.CustomerId, or.ProductId, or.Quantity).Scan(&or.Id)
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
	if err != nil {
		log.Println(err)
		return nil, err
	}
<<<<<<< HEAD
	return &pb.Order{
		Id:                req.Id,
		UserId:            req.UserId,
		ProductIds:        req.ProductIds,
		ProductsDetails:   req.ProductsDetails,
		TotalPayment:      req.TotalPayment,
		SettlementStatus:  req.SettlementStatus,
		IsPrivate:         req.IsPrivate,
		IsPrivateApproved: req.IsPrivateApproved,
		PrivateSellerId:   req.PrivateSellerId,
		CreatedAt:         req.CreatedAt,
		UpdatedAt:         req.UpdatedAt,
	}, nil
}

func (o *OrderQueryImpl) GetOrder(ctx context.Context, req *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	query := `SELECT * from orders WHERE user_id=$1`
	rows, err := o.db.Query(ctx, query, req.CustomerId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no orders found for customer ID %s", req.CustomerId)
		}
		log.Printf("Error querying orders: %v", err)
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
=======
	return &pb.CreateOrderResponse{Id: or.Id, Status: true}, nil
}

func (o *OrderQueryImpl) GetOrder(ctx context.Context, of *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	query := `SELECT * from orders WHERE customer_id=$1`
	rows, err := o.Db.Query(ctx, query, of.CustomerId)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("no records found")
	} else if err != nil {
		return nil, err
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
	}
	var orders []*pb.Order
	for rows.Next() {
		var order pb.Order
<<<<<<< HEAD
		err = rows.Scan(&order.Id, &order.UserId, &order.ProductIds, &order.ProductsDetails, &order.TotalPayment, &order.SettlementStatus, &order.SettlementStatus, &order.IsPrivate, &order.IsPrivateApproved, &order.PrivateSellerId, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return &pb.GetOrderResponse{Orders: orders}, nil
}
func (o *OrderQueryImpl) GetOrders(ctx context.Context, req *pb.GetOrdersRequest, stream pb.OrderService_GetOrdersServer) error {
	query := `SELECT id, user_id, product_ids, products_details, total_payment, settlement_status, is_private, is_private_approved, private_seller_id, created_at, updated_at FROM orders LIMIT $1 OFFSET $2`
	rows, err := o.db.Query(ctx, query, req.Pagination.GetLimit(), req.Pagination.GetOffset())
=======
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
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		order := &pb.Order{}
<<<<<<< HEAD
		if err := rows.Scan(&order.Id, &order.UserId, &order.ProductIds, &order.ProductsDetails, &order.TotalPayment, &order.SettlementStatus, &order.IsPrivate, &order.IsPrivateApproved, &order.PrivateSellerId, &order.CreatedAt, &order.UpdatedAt); err != nil {
=======

		if err := rows.Scan(&order.Id, &order.CustomerId, &order.ProductId, &order.Quantity); err != nil {
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
			return fmt.Errorf("scan error: %v", err)
		}
		response := &pb.GetOrderResponse{
			Orders: []*pb.Order{order},
		}
<<<<<<< HEAD
		if err := stream.Send(response); err != nil {
			return fmt.Errorf("stream send error: %v", err)
=======
		err := stream.Send(response)
		if err != nil {
			return err
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %v", err)
	}
	return nil
}
<<<<<<< HEAD

func (o *OrderQueryImpl) UpdateOrder(ctx context.Context, tx pgx.Tx, order *pb.Order) (*pb.Order, error) {
	query := `UPDATE orders SET user_id = $1, product_ids = $2, products_details = $3, total_payment = $4, settlement_status = $5, is_private = $6, is_private_approved = $7, private_seller_id = $8, updated_at = $9 WHERE id = $10 RETURNING *`
	var updatedOrder pb.Order
	err := tx.QueryRow(ctx, query, order.UserId, order.ProductIds, order.ProductsDetails, order.TotalPayment, order.SettlementStatus, order.IsPrivate, order.IsPrivateApproved, order.PrivateSellerId, order.UpdatedAt, order.Id).Scan(
		&updatedOrder.Id, &updatedOrder.UserId, &updatedOrder.ProductIds, &updatedOrder.ProductsDetails, &updatedOrder.TotalPayment, &updatedOrder.SettlementStatus, &updatedOrder.IsPrivate, &updatedOrder.IsPrivateApproved, &updatedOrder.PrivateSellerId, &updatedOrder.CreatedAt, &updatedOrder.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return &updatedOrder, nil
}

func (o *OrderQueryImpl) SendOrder(ctx context.Context, tx pgx.Tx, req *pb.SendOrderRequest) error {
	query := `UPDATE orders SET settlement_status = 'requires_payment_method', updated_at = $2 WHERE id = $1 RETURNING id, user_id, product_ids, products_details, total_payment, settlement_status, is_private, is_private_approved, private_seller_id, created_at, updated_at`
	var order pb.Order
	err := tx.QueryRow(ctx, query, req.OrderId, time.Now()).Scan(
		&order.Id, &order.UserId, &order.ProductIds, &order.ProductsDetails,
		&order.TotalPayment, &order.SettlementStatus, &order.IsPrivate,
		&order.IsPrivateApproved, &order.PrivateSellerId, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to send order: %w", err)
	}
	return nil
}
=======
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
