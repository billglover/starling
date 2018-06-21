package starling

import (
	"context"
	"net/http"
	"path"
)

// LocalPayment represents a local payment
type LocalPayment struct {
	Payment               PaymentAmount `json:"payment"`
	DestinationAccountUID string        `json:"destinationAccountUid"`
	Reference             string        `json:"reference"`
}

// ScheduledPayment represents a scheduled payment
type ScheduledPayment struct {
	LocalPayment
	Schedule RecurrenceRule `json:"recurrenceRule"`
}

// PaymentOrder is a single PaymentOrder
type PaymentOrder struct {
	UID                        string         `json:"paymentOrderId"`
	Currency                   string         `json:"currency"`
	Amount                     float64        `json:"amount"`
	Reference                  string         `json:"reference"`
	ReceivingContactAccountUID string         `json:"receivingContactAccountId"`
	RecipientName              string         `json:"recipientName"`
	Immediate                  bool           `json:"immediate"`
	RecurrenceRule             RecurrenceRule `json:"recurrenceRule"`
	StartDate                  string         `json:"startDate"`
	NextDate                   string         `json:"nextDate"`
	CancelledAt                string         `json:"cancelledAt"`
	PaymentType                string         `json:"paymentType"`
	MandateUID                 string         `json:"mandateId"`
}

// PaymentOrders is a list of PaymentOrders
type paymentOrders struct {
	PaymentOrders []PaymentOrder `json:"paymentOrders"`
}

// HALPaymentOrders is a HAL wrapper around the Transactions type.
type halPaymentOrders struct {
	Embedded *paymentOrders `json:"_embedded"`
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

// CreateScheduledPayment creates a scheduled payment. It returns the UID for the scheduled payment.
func (c *Client) CreateScheduledPayment(ctx context.Context, p ScheduledPayment) (string, *http.Response, error) {
	req, err := c.NewRequest("POST", "/api/v1/payments/scheduled", p)
	if err != nil {
		return "", nil, err
	}

	resp, err := c.Do(ctx, req, nil)
	if err != nil {
		return "", resp, err
	}

	loc := resp.Header.Get("Location")
	uid := path.Base(loc)
	return uid, resp, err
}

// ScheduledPayments retrieves a list of all the payment orders on the customer account. These may be
// orders for previous immediate payments or scheduled payment orders for future or on-going payments.
func (c *Client) ScheduledPayments(ctx context.Context) ([]PaymentOrder, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/payments/scheduled", nil)

	if err != nil {
		return nil, nil, err
	}

	hPO := new(halPaymentOrders)
	resp, err := c.Do(ctx, req, &hPO)

	if hPO == nil {
		return nil, resp, err
	}

	if hPO.Embedded == nil {
		return nil, resp, err
	}

	return hPO.Embedded.PaymentOrders, resp, err
}
