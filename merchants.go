package starling

import (
	"context"
	"net/http"
)

// Merchant represents details of a merchant
type Merchant struct {
	UID             string `json:"merchantUid"`
	Name            string `json:"name"`
	Website         string `json:"website"`
	PhoneNumber     string `json:"phoneNumber"`
	TwitterUsername string `json:"twitterUsername"`
}

// MerchantLocation represents details of a merchant location
type MerchantLocation struct {
	UID                            string  `json:"merchantLocationUid"`
	MerchantUID                    string  `json:"merchantUid"`
	Merchant                       HALLink `json:"merchant"`
	MerchantName                   string  `json:"merchantName"`
	LocationName                   string  `json:"locationName"`
	Address                        string  `json:"address"`
	PhoneNumber                    string  `json:"phoneNUmber"`
	GooglePlaceID                  string  `json:"googlePlaceId"`
	MastercardMerchantCategoryCode int32   `json:"mastercardMerchantCategoryCode"`
}

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
