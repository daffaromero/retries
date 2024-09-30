package service

import (
	"context"
	"strings"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/common/utils/logger"
	"github.com/daffaromero/retries/services/product-service/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryService interface {
	CreateCategory(context.Context, *pb.Category, string, string) (*pb.Category, error)
	GetCategoryById(context.Context, *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error)
	GetCategories(context.Context, *pb.GetCategoryFilter, pb.ProductService_GetCategoriesServer) error
	UpdateCategory(context.Context, *pb.Category, string, string) (*pb.Category, error)
	DeleteCategory(context.Context, *pb.GetCategoryFilter) (*pb.DeleteCategoryResponse, error)
}

type CategoryServiceImpl struct {
	catRepo repository.CategoryRepository
	logger  *logger.Log
}

func NewCategoryService(catRepo repository.CategoryRepository, logger *logger.Log) CategoryService {
	return &CategoryServiceImpl{
		catRepo: catRepo,
		logger:  logger,
	}
}

func (c *CategoryServiceImpl) CreateCategory(ctx context.Context, cat *pb.Category, name, desc string) (*pb.Category, error) {
	now := timestamppb.Now()
	cat.Id = uuid.New().String()
	cat.Name = name
	cat.Description = desc
	cat.CreatedAt = now
	cat.UpdatedAt = now

	res, err := c.catRepo.CreateCategory(ctx, cat)
	if err != nil {
		c.logger.CustomError("Category creation failed", err)
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Category already exists.")
		}
		return nil, err
	}
	return res, nil
}

func (c *CategoryServiceImpl) GetCategoryById(ctx context.Context, filter *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error) {
	res, err := c.catRepo.GetCategoryById(ctx, filter)
	if err != nil {
		c.logger.CustomError("Failed to get category by ID", err)
		return nil, err
	}
	return res, nil
}

func (c *CategoryServiceImpl) GetCategories(ctx context.Context, filter *pb.GetCategoryFilter, sv pb.ProductService_GetCategoriesServer) error {
	if err := c.catRepo.GetCategories(ctx, filter, sv); err != nil {
		c.logger.CustomError("Failed to get categories", err)
		return err
	}
	return nil
}

func (c *CategoryServiceImpl) UpdateCategory(ctx context.Context, cat *pb.Category, name, desc string) (*pb.Category, error) {
	cat.Name = name
	cat.Description = desc
	cat.UpdatedAt = timestamppb.Now()

	res, err := c.catRepo.UpdateCategory(ctx, cat)
	if err != nil {
		c.logger.CustomError("Failed to update category", err)
		return nil, err
	}
	return res, nil
}

func (c *CategoryServiceImpl) DeleteCategory(ctx context.Context, filter *pb.GetCategoryFilter) (*pb.DeleteCategoryResponse, error) {
	res, err := c.catRepo.DeleteCategory(ctx, filter)
	if err != nil {
		c.logger.CustomError("Failed to delete category", err)
		return nil, err
	}
	return res, nil
}
