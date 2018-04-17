package starling

import (
	"context"
	"net/http"
)

// DirectDebitMandate represents a single mandate
type DirectDebitMandate struct {
	UID            string `json:"uid"`
	Reference      string `json:"reference"`
	Status         string `json:"status"`
	Source         string `json:"source"`
	Created        string `json:"created"`
	Cancelled      string `json:"cancelled"`
	OriginatorName string `json:"originatorName"`
	OriginatorUID  string `json:"originatorUid"`
}

// DirectDebitMandates represents a list of mandates
type directDebitMandates struct {
	Mandates []DirectDebitMandate `json:"mandates"`
}

// HALContacts is a HAL wrapper around the DirectDebitMandates type.
type halDirectDebitMandates struct {
	Links    struct{}             `json:"_links"`
	Embedded *directDebitMandates `json:"_embedded"`
}

// DirectDebitMandates returns the DirectDebitMandates for the current customer.
func (c *Client) DirectDebitMandates(ctx context.Context) (*[]DirectDebitMandate, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/direct-debit/mandates", nil)
	if err != nil {
		return nil, nil, err
	}

	var halResp *halDirectDebitMandates
	var mandates *directDebitMandates
	resp, err := c.Do(ctx, req, &halResp)
	if err != nil {
		return &mandates.Mandates, resp, err
	}

	if halResp.Embedded != nil {
		mandates = halResp.Embedded
	}

	return &mandates.Mandates, resp, nil
}
