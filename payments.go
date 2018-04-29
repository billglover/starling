package starling

import (
	"context"
	"net/http"
)

// LocalPayment represents a local payment
type LocalPayment struct {
	Payment               PaymentAmount `json:"payment"`
	DestinationAccountUID string        `json:"destinationAccountUid"`
	Reference             string        `json:"reference"`
}

// PaymentAmount represents the currency and amount of a payment
type PaymentAmount struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// MakeLocalPayment creates a local payment.
func (c *Client) MakeLocalPayment(ctx context.Context, p LocalPayment) (*http.Response, error) {
	req, err := c.NewRequest("POST", "/api/v1/payments/local", p)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(ctx, req, nil)
	return resp, err
}
