package controller

import (
	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/product-service/config"
	"github.com/daffaromero/retries/services/product-service/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type CategoryController interface {
	Route(*fiber.App)
	CreateCategory(fiber.Ctx) error
	GetCategoryById(fiber.Ctx) error
	GetCategories(fiber.Ctx) error
	UpdateCategory(fiber.Ctx) error
	DeleteCategory(fiber.Ctx) error
}

type CategoryControllerImpl struct {
	validate        *validator.Validate
	categoryService service.CategoryService
}

func NewCategoryController(val *validator.Validate, catServ service.CategoryService) CategoryController {
	return &CategoryControllerImpl{
		validate:        val,
		categoryService: catServ,
	}
}

func (c *CategoryControllerImpl) Route(app *fiber.App) {
	api := app.Group(config.EndpointPrefix)
	api.Post("/category/new", c.CreateCategory)
}

func (c *CategoryControllerImpl) CreateCategory(ctx fiber.Ctx) error {
	var req *pb.Category
	err := ctx.Bind().Body(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	cat, err := c.categoryService.CreateCategory(ctx.Context(), req, req.Name, req.Description)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	res := &pb.Category{
		Id:          cat.Id,
		Name:        cat.Name,
		Description: cat.Description,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.DeletedAt,
		DeletedAt:   cat.DeletedAt,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *CategoryControllerImpl) GetCategoryById(ctx fiber.Ctx) error {
	var req pb.GetCategoryFilter
	req.Id = ctx.Query("id")
	if req.Id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id not provided"})
	}
	cat, err := c.categoryService.GetCategoryById(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get category by id"})
	}
	res := &pb.GetCategoryResponse{
		Categories: cat.Categories,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}
