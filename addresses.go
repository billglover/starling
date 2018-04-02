package starling

import (
	"context"
	"net/http"
)

// Address is the physical address of the customer
type Address struct {
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	Country       string `json:"country"`
	Postcode      string `json:"postcode"`
}

// Addresses are the current and previous physical addresses
type Addresses struct {
	Current  Address   `json:"current"`
	Previous []Address `json:"previous"`
}

// GetAddresses returns the the customer details for the current customer.
func (c *Client) GetAddresses(ctx context.Context) (*Addresses, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/addresses", nil)
	if err != nil {
		return nil, nil, err
	}

	var addr *Addresses
	resp, err := c.Do(ctx, req, &addr)
	if err != nil {
		return addr, resp, err
	}

	return addr, resp, nil
}
