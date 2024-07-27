package main

import (
	"log"
	"net"

	handlers "github.com/daffaromero/retries/services/purchases/handlers/purchases"
	"github.com/daffaromero/retries/services/purchases/service"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
}

func NewgRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

func (gs *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", gs.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRPCServer := grpc.NewServer()

	// registering gRPC services
	purchaseService := service.NewPurchaseService()
	handlers.NewgRPCPurchaseService(gRPCServer, purchaseService)

	log.Println("Starting gRPC server on", gs.addr)

	return gRPCServer.Serve(lis)
}
