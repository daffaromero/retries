package main

// import (
// 	"log"
// 	"net/http"

// 	handlers "github.com/daffaromero/retries/services/order-service/handlers/purchases"
// 	"github.com/daffaromero/retries/services/order-service/service"
// )

// type httpServer struct {
// 	addr string
// }

// func NewHTTPServer(addr string) *httpServer {
// 	return &httpServer{addr: addr}
// }

// func (s *httpServer) Run() error {
// 	router := http.NewServeMux()

// 	purchaseService := service.NewPurchaseService()
// 	purchaseHandler := handlers.NewHTTPPurchaseHandler(purchaseService)
// 	purchaseHandler.RegisterRouter(router)

// 	log.Println("Starting server on", s.addr)

// 	return http.ListenAndServe(s.addr, router)
// }
