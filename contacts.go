package starling

import (
	"context"
	"net/http"
)

// Contact represents the details of a payee
type Contact struct {
	UID  string `json:"id"`
	Name string `json:"name"`
}

// Contacts are a list of payees
type Contacts struct {
	Contacts []Contact
}

// HALContacts is a HAL wrapper around the Contacts type.
type HALContacts struct {
	Links    struct{}  `json:"_links"`
	Embedded *Contacts `json:"_embedded"`
}

// GetContacts returns the contacts for the current customer.
func (c *Client) GetContacts(ctx context.Context) (*Contacts, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/contacts", nil)
	if err != nil {
		return nil, nil, err
	}

	var halResp *HALContacts
	var contacts *Contacts
	resp, err := c.Do(ctx, req, &halResp)
	if err != nil {
		return contacts, resp, err
	}

	if halResp.Embedded != nil {
		contacts = halResp.Embedded
	}

	return contacts, resp, nil
}

// GetContact returns an individual contact for the current customer.
func (c *Client) GetContact(ctx context.Context, uid string) (*Contact, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/contacts/"+uid, nil)
	if err != nil {
		return nil, nil, err
	}

	var contact *Contact
	resp, err := c.Do(ctx, req, &contact)
	return contact, resp, nil
}

// DeleteContact deletes an individual contact for the current customer. It returns http.StatusNoContent
// on success. No payload is returned.
func (c *Client) DeleteContact(ctx context.Context, uid string) (*http.Response, error) {
	req, err := c.NewRequest("DELETE", "/api/v1/contacts/"+uid, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(ctx, req, nil)
	return resp, nil
}
