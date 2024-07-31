package handlers

// import (
// 	"context"

// 	"github.com/daffaromero/retries/services/common/genproto/purchases"
// 	"github.com/daffaromero/retries/services/order-service/types"
// 	"google.golang.org/grpc"
// )

// type PurchasegRPCHandler struct {
// 	purchaseService types.PurchaseService
// 	purchases.UnimplementedPurchaseServiceServer
// }

// func NewgRPCPurchaseService(gs *grpc.Server, ps types.PurchaseService) {
// 	gRPCHandler := &PurchasegRPCHandler{
// 		purchaseService: ps,
// 	}

// 	// register the PurchaseServiceServer
// 	purchases.RegisterPurchaseServiceServer(gs, gRPCHandler)
// }

// func (h *PurchasegRPCHandler) GetOrders(ctx context.Context, req *purchases.GetOrdersRequest) (*purchases.GetOrderResponse, error) {
// 	p := h.purchaseService.GetOrders(ctx)
// 	res := &purchases.GetOrderResponse{
// 		Orders: p,
// 	}

// 	return res, nil
// }

// func (h *PurchasegRPCHandler) CreateOrder(ctx context.Context, req *purchases.CreateOrderRequest) (*purchases.CreateOrderResponse, error) {
// 	purchase := &purchases.Order{
// 		OrderId:    5,
// 		CustomerId: 16,
// 		ProductId:  33,
// 		Quantity:   44,
// 	}

// 	err := h.purchaseService.CreateOrder(ctx, purchase)
// 	if err != nil {
// 		return nil, err
// 	}

// 	res := &purchases.CreateOrderResponse{
// 		Status: "success",
// 	}

// 	return res, nil
// }
