package processor

import (
	"fmt"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/common/utils"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

var gatewayAddress = utils.GetEnv("PAYMENT_GATEWAY_ADDR")

type Stripe struct{}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(o *pb.SendOrderRequest) (string, error) {
	gatewaySuccessURL := fmt.Sprintf("%s/payment/success.html?&orderID=%s", gatewayAddress, o.OrderId)
	gatewayCancelURL := fmt.Sprintf("%s/payment/cancel.html", gatewayAddress)

	params := &stripe.CheckoutSessionParams{
		Metadata: map[string]string{
			"order_id": o.OrderId,
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
		CancelURL:  stripe.String(gatewayCancelURL),
	}

	res, err := session.New(params)
	if err != nil {
		return "", err
	}
	return res.URL, nil
}
