package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	api.Get("/category/:id", c.GetCategoryById)
	api.Get("/category", c.GetCategories)
	api.Put("/category/:id", c.UpdateCategory)
	api.Delete("/category/:id", c.DeleteCategory)
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

func (c *CategoryControllerImpl) GetCategories(ctx fiber.Ctx) error {
	var req pb.GetCategoryFilter
	err := ctx.Bind().Body(&req)
	if err != nil {
		return fmt.Errorf("error binding request - %s", err)
	}
	count, _ := strconv.Atoi(ctx.Query("count"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	req.Pagination.Count = int32(count)
	req.Pagination.Limit = int32(limit)
	req.Pagination.Page = int32(page)

	grpcServer := NewRestCategoryServer()

	go func() {
		if err := c.categoryService.GetCategories(ctx.Context(), &req, grpcServer); err != nil {
			log.Printf("error getting categories: %v", err)
		}
		close(grpcServer.results)
	}()

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		enc := json.NewEncoder(w)
		for {
			res, err := grpcServer.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("error receiving categories response: %v", err)
				ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get categoris"})
				return
			}
			if err := enc.Encode(res); err != nil {
				log.Printf("error encoding categories response: %v", err)
				ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to encode categories"})
				return
			}
			w.Flush()
		}
	})
	return nil
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
