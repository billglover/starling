package starling

import (
	"context"
	"net/http"
)

// Identity is the identity of the current user
type Identity struct {
	UID              string   `json:"customerUid"`
	ExpiresAt        string   `json:"expiresAt"`
	Authenticated    bool     `json:"authenticated"`
	ExpiresInSeconds int64    `json:"expiresInSeconds"`
	Scopes           []string `json:"scopes"`
}

// CurrentUser returns the identity of the current user.
func (c *Client) CurrentUser(ctx context.Context) (*Identity, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/me", nil)
	if err != nil {
		return nil, nil, err
	}

	var ident *Identity
	resp, err := c.Do(ctx, req, &ident)
	if err != nil {
		return ident, resp, err
	}

	return ident, resp, nil
}
