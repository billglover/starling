package starling

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.baseURL.String(), defaultURL; got != want {
		t.Errorf("NewClient baseURL is %v, want %v", got, want)
	}
	if got, want := c.userAgent, userAgent; got != want {
		t.Errorf("NewClient userAgent is %v, want %v", got, want)
	}
}

// TestNewRequest confirms that NewRequest returns an API request with the
// correct URL, a correctly encoded body and the correct User-Agent and
// Content-Type headers set.
func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultURL+"foo"
	inBody, outBody := &TopUpRequest{Amount: CurrencyAndAmount{Currency: "GBP", MinorUnits: 1973}}, `{"amount":{"currency":"GBP","minorUnits":1973}}`+"\n"
	req, err := c.NewRequest("GET", inURL, inBody)
	if err != nil {
		t.Fatalf("NewRequest(%q) resulted in an error: %v", inURL, err)
	}

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	// test that default user-agent is set
	if got, want := req.Header.Get("User-Agent"), c.userAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}

	// test that default content-type is set
	if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("NewRequest() Content-Type is %v, want %v", got, want)
	}
}

// TestNewRequest_invalidJSON confirms that NewRequest returns an error
// if asked to encode a type that results in invalid JSON.
func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.NewRequest("GET", ".", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}

}

// TestNewRequest_badURL confirms that NewRequest returns an error
// if passed a URL containing invalid characters.
func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("GET", ":", nil)
	if err == nil {
		t.Error("expected error to be returned")
	}
}

// TestNewRequest_badBasePath confirms that NewRequest returns an error
// if called on a client that does not have a trailing slash for the
// base path.
func TestNewRequest_badBasePath(t *testing.T) {
	c := NewClient(nil)
	u, _ := url.Parse("http://test.local")
	c.baseURL = u
	_, err := c.NewRequest("GET", "/", nil)
	if err == nil {
		t.Error("expected error to be returned")
	}
}

// TestNewRequest_badMethod confirms that NewRequest returns an error
// if called with an invalid method.
func TestNewRequest_badMethod(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("\n", "/", nil)
	if err == nil {
		t.Error("expected error to be returned")
	}
}

// TestNewRequest_emptyBody confirms that NewRequest returns an API request with the
// correct URL, an empty body and the correct User-Agent and Content-Type headers set.
func TestNewRequest_emptyBody(t *testing.T) {
	c := NewClient(nil)
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(foo)
	client.Do(context.Background(), req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	// client is the GitHub client being tested and is
	// configured to use test server.
	c := NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	c.baseURL = url

	return c, mux, server.URL, server.Close
}
