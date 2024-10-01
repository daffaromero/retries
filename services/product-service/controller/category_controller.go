package controller

import (
	"fmt"
	"strconv"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/product-service/config"
	"github.com/daffaromero/retries/services/product-service/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type CategoryController interface {
	Route(*fiber.App)
	CreateCategory(fiber.Ctx) error
	GetCategoryByID(fiber.Ctx) error
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
	api.Post("/new", c.CreateCategory)
	api.Get("/:id", c.GetCategoryByID)
	api.Get("/", c.GetCategories)
	api.Put("/:id", c.UpdateCategory)
	api.Delete("/:id", c.DeleteCategory)
}

func (c *CategoryControllerImpl) CreateCategory(ctx fiber.Ctx) error {
	var req pb.Category
	err := ctx.Bind().Body(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.validate.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Id != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id must be provided"})
	}

	res, err := c.categoryService.CreateCategory(ctx.Context(), &req, req.Name, req.Description)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *CategoryControllerImpl) GetCategoryByID(ctx fiber.Ctx) error {
	var req pb.GetCategoryFilter
	req.Id = ctx.Query("id")
	if req.Id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id not provided"})
	}
	cat, err := c.categoryService.GetCategoryByID(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get category by id"})
	}
	res := &pb.GetCategoryResponse{
		Categories: cat.Categories,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *CategoryControllerImpl) GetCategories(ctx fiber.Ctx) error {
	var fil pb.GetCategoryFilter
	err := ctx.Bind().Body(&fil)
	if err != nil {
		return fmt.Errorf("error binding request - %s", err)
	}
	offset, _ := strconv.Atoi(ctx.Query("offset"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	fil.Pagination.Offset = int32(offset)
	fil.Pagination.Limit = int32(limit)
	fil.Pagination.Page = int32(page)

	categories, err := c.categoryService.GetCategories(ctx.Context(), &fil)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(categories)
}

func (c *CategoryControllerImpl) UpdateCategory(ctx fiber.Ctx) error {
	var req *pb.Category
	req.Id = ctx.Query("id")
	if req.Id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id not provided"})
	}
	err := ctx.Bind().Body(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}
	cat, err := c.categoryService.UpdateCategory(ctx.Context(), req, req.Name, req.Description)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	res := &pb.Category{
		Id:          cat.Id,
		Name:        cat.Name,
		Description: cat.Description,
		UpdatedAt:   cat.UpdatedAt,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *CategoryControllerImpl) DeleteCategory(ctx fiber.Ctx) error {
	var req *pb.GetCategoryFilter
	req.Id = ctx.Query("id")
	if req.Id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id not provided"})
	}
	err := ctx.Bind().Body(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}
	cat, err := c.categoryService.DeleteCategory(ctx.Context(), req)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	res := &pb.DeleteCategoryResponse{
		Status: cat.Status,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}
