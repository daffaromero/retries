package service

import (
	"context"

	pb "github.com/daffaromero/retries/services/common/genproto/orders"
	"github.com/daffaromero/retries/services/common/utils/logger"
	"github.com/daffaromero/retries/services/order-service/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, string, string, int32) (*pb.CreateOrderResponse, error)
	GetOrder(context.Context, *pb.GetOrderFilter, string) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest, pb.OrderService_GetOrdersServer) error
}

type orderService struct {
	ordRepo repository.OrderRepository
	logger  *logger.Log
}

func NewOrderService(ordRepo repository.OrderRepository, logger *logger.Log) OrderService {
	return &orderService{
		ordRepo: ordRepo,
		logger:  logger,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, order *pb.CreateOrderRequest, customerId string, productId string, quantity int32) (*pb.CreateOrderResponse, error) {
	order.Id = uuid.New().String()
	order.CustomerId = customerId
	order.ProductId = productId
	order.Quantity = quantity

	res, err := o.ordRepo.CreateOrder(ctx, order)
	if err != nil {
		o.logger.CustomError("Order creation failed", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create order, please try again.")
	}

	return res, nil
}

func (o *orderService) GetOrder(ctx context.Context, filter *pb.GetOrderFilter, customerId string) (*pb.GetOrderResponse, error) {
	
}