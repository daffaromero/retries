package service

import (
	"context"
	"log"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/common/utils/logger"
	"github.com/daffaromero/retries/services/order-service/repository"
	"github.com/daffaromero/retries/services/product-service/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/stripe/stripe-go/v79"
)

type OrderService interface {
	CreateOrder(c context.Context, ord *pb.Order, id string, name string, phone string, email string) (*pb.Order, error)
	GetOrder(context.Context, *pb.GetOrderFilter) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest) (*pb.GetOrderResponse, error)
	SendOrder(context.Context, *pb.SendOrderRequest) (*stripe.PaymentLink, error)
}

type orderService struct {
	ordRepo repository.OrderRepository
	prodRepo repository.ProductRepository
	logger  *logger.Log
}

func NewOrderService(ordRepo repository.OrderRepository, logger *logger.Log) OrderService {
	return &orderService{
		ordRepo: ordRepo,
		logger:  logger,
	}
}

func (o *orderService) CreateOrder(c context.Context, ord *pb.Order) (*pb.Order, error) {
	var prDet []*pb.ProductDetails
	var prIds []string
	var total int

	if len(ord.ProductIds) == 0 {
		for _, p := range ord.ProductsDetails {
			product, err := o.productRepo
		}
	}

	res, err := o.ordRepo.CreateOrder(ctx, order)
	if err != nil {
		o.logger.CustomError("Order creation failed", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create order, please try again.")
	}
	return res, nil
}

func (o *orderService) GetOrder(ctx context.Context, filter *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	res, err := o.ordRepo.GetOrder(ctx, filter)
	if err != nil {
		o.logger.CustomError("Failed to get order by ID", err)
		return nil, err
	}
	return res, nil
}

func (o *orderService) GetAllOrders(ctx context.Context, req *pb.GetOrdersRequest, sm pb.OrderService_GetOrdersServer) error {
	log.Print(ctx)
	err := o.ordRepo.GetAllOrders(ctx, req, sm)
	if err != nil {
		o.logger.CustomError("failed to get all orders", err)
		return err
	}
	return nil
}
