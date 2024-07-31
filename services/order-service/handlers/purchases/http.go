package handlers

// import (
// 	"net/http"

// 	"github.com/daffaromero/retries/services/common/genproto/purchases"
// 	"github.com/daffaromero/retries/services/common/utils"
// 	"github.com/daffaromero/retries/services/order-service/types"
// )

// type PurchaseHTTPHandler struct {
// 	purchaseService types.PurchaseService
// }

// func NewHTTPPurchaseHandler(ps types.PurchaseService) *PurchaseHTTPHandler {
// 	handler := &PurchaseHTTPHandler{
// 		purchaseService: ps,
// 	}

// 	return handler
// }

// func (h *PurchaseHTTPHandler) RegisterRouter(router *http.ServeMux) {
// 	router.HandleFunc("POST /purchase", h.CreateOrder)
// }

// func (h *PurchaseHTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
// 	var req purchases.CreateOrderRequest
// 	err := utils.ParseJSON(r, &req)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	order := &purchases.Order{
// 		OrderId:    5,
// 		CustomerId: req.GetCustomerId(),
// 		ProductId:  req.GetProductId(),
// 		Quantity:   req.GetQuantity(),
// 	}

// 	err = h.purchaseService.CreateOrder(r.Context(), order)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	res := &purchases.CreateOrderResponse{Status: "success"}
// 	utils.WriteJSON(w, http.StatusOK, res)
// }
