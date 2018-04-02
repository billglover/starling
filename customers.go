package starling

import (
	"context"
	"net/http"
)

// Customer represents the personal details of a customer
type Customer struct {
	UID               string `json:"customerUid"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	DateOfBirth       string `json:"dateOfBirth"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	AccountHolderType string `json:"accountHolderType"`
}

// GetCustomer returns the the customer details for the current customer.
func (c *Client) GetCustomer(ctx context.Context) (*Customer, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/customers", nil)
	if err != nil {
		return nil, nil, err
	}

	var cust *Customer
	resp, err := c.Do(ctx, req, &cust)
	if err != nil {
		return cust, resp, err
	}

	return cust, resp, nil
}
