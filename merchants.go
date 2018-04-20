package starling

import (
	"context"
	"net/http"
)

// /api/v1/merchants/c052f76f-e919-427d-85fc-f46a75a3ff26
// /api/v1/merchants/c052f76f-e919-427d-85fc-f46a75a3ff26/locations/371c62bc-dcfc-4799-8b23-b070626772f7

// Merchant returns an individual merchant based on the UID.
func (c *Client) Merchant(ctx context.Context, uid string) (*Merchant, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/merchants/"+uid, nil)
	if err != nil {
		return nil, nil, err
	}

	mer := new(Merchant)
	resp, err := c.Do(ctx, req, mer)
	return mer, resp, err
}

// MerchantLocation returns an individual merchant location based on the merchant UID and location UID.
func (c *Client) MerchantLocation(ctx context.Context, mUID, lUID string) (*MerchantLocation, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/merchants/"+mUID+"/locations/"+lUID, nil)
	if err != nil {
		return nil, nil, err
	}

	merLoc := new(MerchantLocation)
	resp, err := c.Do(ctx, req, merLoc)
	return merLoc, resp, err
}
