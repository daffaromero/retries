package service

import (
	"context"
	"fmt"
	"math"

	"github.com/daffaromero/retries/services/common/discovery"
	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/common/utils/logger"
	"github.com/daffaromero/retries/services/order-service/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/stripe/stripe-go/v79"
)

type OrderService interface {
	CreateOrder(c context.Context, ord *pb.Order, id string, name string, phone string, email string) (*pb.Order, error)
	GetOrderDetails(context.Context, *pb.GetOrderFilter) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest) (*pb.GetOrderResponse, error)
	SendOrder(context.Context, *pb.SendOrderRequest) (*stripe.PaymentLink, error)
}

type orderService struct {
	client   pb.ProductServiceClient
	registry discovery.Registry
	ordRepo  repository.OrderRepository
	logger   *logger.Log
}

func NewOrderService(ctx context.Context, registry discovery.Registry, ordRepo repository.OrderRepository, logger *logger.Log) (*orderService, error) {
	conn, err := discovery.ConnectToService(ctx, "product-service-grpc", registry)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to product service", err))
		return nil, err
	}
	return &orderService{
		client:   pb.NewProductServiceClient(conn),
		registry: registry,
		ordRepo:  ordRepo,
		logger:   logger,
	}, nil
}

func (o *orderService) CreateOrder(c context.Context, ord *pb.Order) (*pb.Order, error) {
	var prDet []*pb.ProductDetails
	var prIds []string
	var total int

	if o.client == nil {
		o.logger.Error("Product service client is not initialized")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create order, please try again.")
	}

	if len(ord.ProductIds) == 0 {
		for _, p := range ord.ProductsDetails {
			product, err := o.client.GetProductByID(c, &pb.GetProductFilter{Id: p.Id})
			if err != nil {
				return nil, fiber.NewError(fiber.StatusNotFound, "No product found with ID %v", p.Id)
			}

			if p.Voucher != "" {
				for _, prod := range product.Products {
					if prod.Price <= 5 {
						return nil, fiber.NewError(fiber.StatusNotFound, "Vouchers can not be applied to products lesser than 5 dollars")
					}
					if p.Voucher != prod.Voucher {
						return nil, fiber.NewError(fiber.StatusNotFound, "Voucher not applicable.")
					}

					p.Price -= (p.Price * p.VoucherDiscount / 100)
				}
			}

			prDet = append(prDet, &pb.ProductDetails{
				Id:         p.Id,
				Name:       p.Name,
				Category:   p.Category,
				SellerName: p.SellerName,
				Price:      p.Price,
			})
			total += int(p.Price)
			prIds = append(prIds, p.Id)
		}

		if total < 5 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Total can not be less than 5 dollars")
		} else {
			total = 0
		}
	}

	if total != 0 {
		if (float64(total) * 2 / 100) > 5 {
			adminFee := float64(total) * 2 / 100
			ceilFee := math.Ceil(adminFee)
			ord.AdminFee = int32(ceilFee)
			ord.GrandTotal = int32(total + int(ceilFee))
		}
	} else {
		ord.AdminFee = 0
		ord.GrandTotal = 0
	}

	res, err := o.ordRepo.CreateOrder(c, ord)
	if err != nil {
		o.logger.CustomError("Order creation failed", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create order, please try again.")
	}
	return res, nil
}

func (o *orderService) GetOrderDetails(ctx context.Context, filter *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	order, err := o.ordRepo.GetOrderDetails(ctx, filter)
	if err != nil {
		o.logger.CustomError("Failed to get order by ID", err)
		return nil, err
	}
	return order, nil
}

func (o *orderService) GetAllOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrderResponse, error) {
	orders, err := o.ordRepo.GetAllOrders(ctx, req)
	if err != nil {
		o.logger.CustomError("failed to get all orders", err)
		return nil, err
	}
	return orders, nil
}

func (o *orderService) SendOrder(ctx context.Context, req *pb.SendOrderRequest) (*stripe.PaymentLink, error) {
	link, err := o.ordRepo.SendOrder(ctx, req)
	if err != nil {
		o.logger.CustomError("Failed to send order", err)
		return nil, err
	}
	return link, nil
}