package service

import (
	"context"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/common/utils/logger"
	"github.com/daffaromero/retries/services/product-service/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService interface {
	CreateProduct(context.Context, *pb.Product, string, string) (*pb.Product, error)
	GetProductByID(context.Context, *pb.GetProductFilter) (*pb.GetProductResponse, error)
	GetAllProducts(context.Context, *pb.GetProductFilter) (*pb.GetProductResponse, error)
	UpdateProduct(context.Context, *pb.Product) (*pb.Product, error)
	ApproveProduct(context.Context, *pb.ApproveProductRequest) (*pb.ApproveProductResponse, error)
}

type productService struct {
	productRepo repository.ProductRepository
	logger      *logger.Log
}

func NewProductService(productRepo repository.ProductRepository, logger *logger.Log) ProductService {
	return &productService{
		productRepo: productRepo,
		logger:      logger,
	}
}

func (p *productService) CreateProduct(c context.Context, product *pb.Product, sellerId, sellerName string) (*pb.Product, error) {
	var varIDs []string
	for _, v := range product.VariantSettings {
		varIDs = append(varIDs, v.Id)
	}
	product.Id = uuid.New().String()
	product.SellerId = sellerId
	product.SellerName = sellerName
	product.VariantIds = varIDs
	product.IsAdminVerified = "pending"
	product.CreatedAt = timestamppb.Now()
	product.UpdatedAt = timestamppb.Now()

	res, err := p.productRepo.CreateProduct(c, product)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (p *productService) GetProductByID(c context.Context, filter *pb.GetProductFilter) (*pb.GetProductResponse, error) {
	res, err := p.productRepo.GetProductByID(c, filter)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (p *productService) GetAllProducts(c context.Context, filter *pb.GetProductFilter) (*pb.GetProductResponse, error) {
	res, err := p.productRepo.GetAllProducts(c, filter)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (p *productService) UpdateProduct(c context.Context, product *pb.Product) (*pb.Product, error) {
	pro, _ := p.productRepo.GetProductByID(c, &pb.GetProductFilter{Id: product.Id})
	if pro.Products[0].Visibility == "active" {
		return nil, fiber.NewError(fiber.StatusForbidden, "product is already active")
	}
	var varIDs []string
	for _, v := range product.VariantSettings {
		varIDs = append(varIDs, v.Id)
	}
	product.VariantIds = varIDs
	product.IsAdminVerified = "pending"
	product.UpdatedAt = timestamppb.Now()

	res, err := p.productRepo.UpdateProduct(c, product)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (p *productService) ApproveProduct(c context.Context, req *pb.ApproveProductRequest) (*pb.ApproveProductResponse, error) {
	res, err := p.productRepo.ApproveProduct(c, req)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return res, nil
}
