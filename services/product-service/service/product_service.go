package service

import (
	"context"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/common/utils/logger"
	"github.com/daffaromero/retries/services/product-service/repository"
	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(context.Context, *pb.Product, string, string) (*pb.Product, error)
	GetProductById(context.Context, *pb.GetProductFilter) (*pb.GetProductResponse, error)
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
	//TODO finish
}
