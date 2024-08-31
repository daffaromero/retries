package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/daffaromero/retries/services/common/utils/logger"
	"github.com/daffaromero/retries/services/order-service/config"
	"github.com/daffaromero/retries/services/order-service/controller"
	"github.com/daffaromero/retries/services/order-service/repository"
	"github.com/daffaromero/retries/services/order-service/repository/query"
	"github.com/daffaromero/retries/services/order-service/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	flog "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

var logs = logger.NewLog("main")

func webServer() error {
	app := fiber.New(fiber.Config{
		StreamRequestBody: true,
	})

	app.Use(requestid.New())
	app.Use(flog.New())

	serverConfig := config.NewServerConfig()
	dbConfig := config.NewPGDatabase()
	store := repository.NewStore(dbConfig)
	validate := validator.New()

	ordQuery := query.NewOrderQueryImpl(dbConfig)
	ordRepo := repository.NewOrderRepository(store, ordQuery)
	ordServ := service.NewOrderService(ordRepo, logs)
	ordCont := controller.NewOrderController(validate, ordServ)

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			allowedOrigins := []string{
				"http://localhost",
				"http://localhost:3000",
			}
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return true
		},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodOptions,
		},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization",
		},
		MaxAge: 0,
	}))

	ordCont.Route(app)

	err := app.Listen(serverConfig.Host)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func main() {
	if err := webServer(); err != nil {
		log.Fatalf("webServer failed: %v", err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
}
