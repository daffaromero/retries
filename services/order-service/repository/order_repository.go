package repository

import (
	"context"
	"fmt"
	"log"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/order-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/v79"
)

type OrderRepository interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	GetOrder(context.Context, *pb.GetOrderFilter) (*pb.GetOrderResponse, error)
	GetAllOrders(context.Context, *pb.GetOrdersRequest, pb.OrderService_GetOrdersServer) error
	SendOrder(context.Context, pgx.Tx, *pb.SendOrderRequest) (*stripe.PaymentIntent, error)
}

type orderRepository struct {
	db       Store
	ordQuery query.OrderQuery
}

func NewOrderRepository(db Store, ordQuery query.OrderQuery) OrderRepository {
	return &orderRepository{db: db, ordQuery: ordQuery}
}

func (o *orderRepository) CreateOrder(ctx context.Context, ge *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var res *pb.CreateOrderResponse
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		re, err := o.ordQuery.CreateOrder(ctx, ge)
		if err != nil {
			return err
		}
		res = re
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderRepository) GetOrder(ctx context.Context, gf *pb.GetOrderFilter) (*pb.GetOrderResponse, error) {
	var res *pb.GetOrderResponse
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		re, err := o.ordQuery.GetOrder(ctx, gf)
		if err != nil {
			return err
		}
		res = re
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderRepository) GetAllOrders(ctx context.Context, req *pb.GetOrdersRequest, sm pb.OrderService_GetOrdersServer) error {
	err := o.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		return o.ordQuery.GetAllOrders(ctx, req, sm)
	})
	if err != nil {
		log.Print("the repo brokey")
		return err
	}
	return nil
}

func (o *orderRepository) SendOrder(ctx context.Context, tx pgx.Tx, req *pb.SendOrderRequest) error {

	// todo fix the whole file lol
	params := &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Metadata: map[string]string{
				"order_id": order.Id,
			},
		},
		Amount:               stripe.Int64(int64(order.TotalPayment)),
		ApplicationFeeAmount: nil,
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		CaptureMethod:              stripe.String(string(stripe.PaymentIntentCaptureMethodAutomatic)),
		ClientSecret:               nil,
		Confirm:                    stripe.Bool(true),
		ConfirmationMethod:         stripe.String(string(stripe.PaymentIntentConfirmationMethodAutomatic)),
		ConfirmationToken:          nil,
		Currency:                   stripe.String("idr"),
		Customer:                   stripe.String(order.UserId),
		Description:                stripe.String("Payment for order " + order.Id),
		Expand:                     []*string{stripe.String("customer")},
		Mandate:                    nil,
		MandateData:                nil,
		Metadata:                   map[string]string{"order_id": order.Id},
		OnBehalfOf:                 nil,
		PaymentMethod:              nil,
		PaymentMethodConfiguration: nil,
		PaymentMethodData:          nil,
		PaymentMethodOptions:       nil,
		PaymentMethodTypes:         []*string{stripe.String("card")},
		RadarOptions:               nil,
		ReceiptEmail:               stripe.String(order.UserId), // Assuming UserId is the email, replace with actual email if available
		ReturnURL:                  stripe.String("https://your-return-url.com"),
		SetupFutureUsage:           stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession)),
		Shipping:                   nil,
		StatementDescriptor:        stripe.String("Your descriptor"),
		StatementDescriptorSuffix:  stripe.String("Order"),
		TransferData:               nil,
		TransferGroup:              nil,
		ErrorOnRequiresAction:      stripe.Bool(false),
		OffSession:                 stripe.Bool(true),
		UseStripeSDK:               stripe.Bool(true),
	}

	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

}
