package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	pb "github.com/daffaromero/retries/services/common/genproto/orders"
	"github.com/daffaromero/retries/services/order-service/config"
	"github.com/daffaromero/retries/services/order-service/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type OrderController interface {
	Route(*fiber.App)
	CreateOrder(fiber.Ctx) error
	GetOrder(fiber.Ctx) error
	GetAllOrders(fiber.Ctx) error
}

type orderController struct {
	validate     *validator.Validate
	orderService service.OrderService
}

func NewOrderController(val *validator.Validate, ordServ service.OrderService) OrderController {
	return &orderController{
		validate:     val,
		orderService: ordServ,
	}
}

func (o *orderController) Route(app *fiber.App) {
	api := app.Group(config.EndpointPrefix)
	api.Post("/new", o.CreateOrder)
	api.Get("/customer", o.GetOrder)
	api.Get("/all", o.GetAllOrders)
}

func (o *orderController) CreateOrder(c fiber.Ctx) error {
	var req *pb.CreateOrderRequest
	err := c.Bind().Body(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	ord, err := o.orderService.CreateOrder(c.Context(), req, req.CustomerId, req.ProductId, req.Quantity)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	res := &pb.CreateOrderResponse{
		Id:     ord.Id,
		Status: true,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (o *orderController) GetOrder(c fiber.Ctx) error {
	var req pb.GetOrderFilter
	req.CustomerId = c.Query("customer_id")
	if req.CustomerId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer_id not provided"})
	}
	ord, err := o.orderService.GetOrder(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get order for customer"})
	}
	res := &pb.GetOrderResponse{
		Orders: ord.Orders,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (o *orderController) GetAllOrders(c fiber.Ctx) error {
	var req pb.GetOrdersRequest
	err := c.Bind().Body(&req)
	if err != nil {
		return fmt.Errorf("error binding request - %s", err)
	}
	req.CustomerId = c.Query("customer_id")
	count, _ := strconv.Atoi(c.Query("count"))
	start, _ := strconv.Atoi(c.Query("start"))
	req.Count = int32(count)
	req.Start = int32(start)

	grpcServer := NewRestOrderServer()

	go func() {
		if err := o.orderService.GetAllOrders(c.Context(), &req, grpcServer); err != nil {
			log.Printf("error getting orders: %v", err)
		}
		close(grpcServer.results)
	}()

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		enc := json.NewEncoder(w)
		for {
			res, err := grpcServer.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("error receiving order response: %v", err)
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get orders"})
				return
			}
			if err := enc.Encode(res); err != nil {
				log.Printf("error encoding order response: %v", err)
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to encode orders"})
				return
			}
			w.Flush()
		}
	})
	return nil
}
