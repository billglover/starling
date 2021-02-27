package starling

import (
	"context"
	"net/http"
)

// Address is the physical address of the customer
type Address struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	Line3   string `json:"line3"`
	PostTown     string `json:"postTown"`
	CountryCode  string `json:"countryCode"`
	PostCode string `json:"postCode"`
}

// AddressHistory are the current and previous physical addresses
type AddressHistory struct {
	Current  Address   `json:"current"`
	Previous []Address `json:"previous"`
}

// AddressHistory returns the the customer details for the current customer.
func (c *Client) AddressHistory(ctx context.Context) (*AddressHistory, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v2/addresses", nil)
	if err != nil {
		return nil, nil, err
	}

	var addr *AddressHistory
	resp, err := c.Do(ctx, req, &addr)
	if err != nil {
		return nil, resp, err
	}

	return addr, resp, nil
}
