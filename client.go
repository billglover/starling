package starling

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	prodURL    = "https://api.starlingbank.com/"
	sandboxURL = "https://api-sandbox.starlingbank.com/"
	defaultURL = prodURL

	userAgent = "go-starling"
)

// Client holds configuration items for the Starling client and provides methods
// that interact with the Starling API.
type Client struct {
	baseURL *url.URL

	userAgent string
	client    *http.Client
}

// NewClient returns a new Starling API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
// Inspiration: https://github.com/google/go-github/blob/master/github/github.go
func NewClient(cc *http.Client) *Client {
	if cc == nil {
		cc = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultURL)

	c := &Client{baseURL: baseURL, userAgent: userAgent, client: cc}
	return c
}

// NewRequest creates an HTTP Request. The client baseURL is checked to confirm that it has a trailing
// slash. A relative URL should be provided without the leading slash. If a non-nil body is provided
// it will be JSON encoded and included in the request.
// Inspiration: https://github.com/google/go-github/blob/master/github/github.go
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if strings.HasSuffix(c.baseURL.Path, "/") == false {
		return nil, fmt.Errorf("client baseURL does not have a trailing slash: %q", c.baseURL)
	}

	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}

// Do sends a request and returns the response. An error is returned if the request cannot
// be sent or if the API returns an error. If a response is received, the body response body
// is decoded and stored in the value pointed to by v.
// Inspiration: https://github.com/google/go-github/blob/master/github/github.go
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)

	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// Anything other than a HTTP 2xx response code is treated as an error. But the structure of error
	// responses differs depending on the API being called. Some APIs return validation errors as part
	// of the standard response. Others respond with a standardised error structure.
	if c := resp.StatusCode; c >= 300 {

		// Try parsing the response using the standard error schema and returning the error.
		var e = ErrorDetail{}
		err := json.Unmarshal(data, &e)
		if err != nil {
			return resp, fmt.Errorf("API returned an error but client was unable to parse the detail: %v", err)
		}

		if e.Message != "" {
			return resp, fmt.Errorf(e.Message)
		}

		// If we haven't been able to parse the standard error schema try parsing the response.
		err = json.Unmarshal(data, v)
		if err != nil {
			return resp, fmt.Errorf("API returned an error but client was unable to parse the detail: %v", err)
		}

		// There isn't much more we can do to determine the cause of the error so return the HTTP status code.
		return resp, fmt.Errorf(resp.Status)
	}

	if v != nil && len(data) != 0 {
		err = json.Unmarshal(data, v)
		if err == io.EOF {
			err = nil
		}
	}

	return resp, err
}
