package starling

import (
	"context"
	"net/http"
)

// Card represents card details
type Card struct {
	UID                 string  `json:"id"`
	NameOnCard          string  `json:"nameOnCard"`
	Type                string  `json:"type"`
	Enabled             bool    `json:"enabled"`
	Cancelled           bool    `json:"cancelled"`
	ActivationRequested bool    `json:"activationRequested"`
	Activated           bool    `json:"activated"`
	DispatchDate        string  `json:"dispatchDate"`
	LastFourDigits      string  `json:"lastFourDigits"`
	Transactions        HALLink `json:"transactions"`
}

// Card returns the the customer details for the current customer.
func (c *Client) Card(ctx context.Context) (*Card, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/cards", nil)
	if err != nil {
		return nil, nil, err
	}

	var card *Card
	resp, err := c.Do(ctx, req, &card)
	if err != nil {
		return card, resp, err
	}

	return card, resp, nil
}
