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

const (
	tick  = "\u2713"
	cross = "\u2717"
)

// TestNewClient confirms that a client can be created with the default baseURL
// and default User-Agent.
func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.baseURL.String(), defaultURL; got != want {
		t.Error("should configure the client to use the default url", cross, got)
	}

	if got, want := c.userAgent, userAgent; got != want {
		t.Error("should configure the client to use the default user-agent", cross, got)
	}
}

// TestNewClientWithOptions confirms that a new client can be created by passing
// in custom ClientOptions.
func TestNewClientWithOptions(t *testing.T) {
	baseURL, _ := url.Parse("https://dummyurl:4000")
	opts := ClientOptions{
		BaseURL: baseURL,
	}
	c := NewClientWithOptions(nil, opts)

	if got, want := c.baseURL.String(), baseURL.String(); got != want {
		t.Error("should configure the client to use the custom url", cross, got)
	}
}

// TestNewRequest confirms that NewRequest returns an API request with the
// correct URL, a correctly encoded body and the correct User-Agent and
// Content-Type headers set.
func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	t.Run("valid request", func(tc *testing.T) {

		inURL, outURL := "/foo", defaultURL+"foo"
		inBody, outBody := &topUpRequest{Amount: Amount{Currency: "GBP", MinorUnits: 1973}}, `{"amount":{"currency":"GBP","minorUnits":1973}}`+"\n"

		req, err := c.NewRequest("GET", inURL, inBody)
		checkNoError(tc, err)

		if got, want := req.URL.String(), outURL; got != want {
			t.Error("should expand relative URLs to absolute URLs", cross, got)
		}

		body, _ := ioutil.ReadAll(req.Body)
		if got, want := string(body), outBody; got != want {
			tc.Error("should encode the request body as JSON", cross, got)
		}

		if got, want := req.Header.Get("User-Agent"), c.userAgent; got != want {
			tc.Error("should pass the correct user agent", cross, got)
		}

		if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
			tc.Error("should set the content-type as application/json", cross, got)
		}
	})

	t.Run("request with invalid JSON", func(tc *testing.T) {
		type T struct{ A map[interface{}]interface{} }
		_, err := c.NewRequest("GET", ".", &T{})
		checkHasError(tc, err)
	})

	t.Run("request with an invalid URL", func(tc *testing.T) {
		_, err := c.NewRequest("GET", ":", nil)
		checkHasError(tc, err)
	})

	t.Run("request with an invalid base path", func(tc *testing.T) {
		u, _ := url.Parse("http://test.local")
		o := ClientOptions{BaseURL: u}
		lc := NewClientWithOptions(nil, o)

		_, err := lc.NewRequest("GET", "/", nil)
		checkHasError(tc, err)
	})

	t.Run("request with an invalid Method", func(tc *testing.T) {
		_, err := c.NewRequest("\n", "/", nil)
		checkHasError(tc, err)
	})

	t.Run("request with an empty body", func(tc *testing.T) {
		req, err := c.NewRequest("GET", ".", nil)
		checkNoError(tc, err)

		if req.Body != nil {
			tc.Error("should return an empty body", cross)
		}
	})

}

// TestDo confirms that Do returns a JSON decoded value when making a request. It
// confirms the correct verb was used and that the decoded response value matches
// the expected result.
func TestDo(t *testing.T) {
	t.Run("successful GET request", func(tc *testing.T) {
		client, mux, _, teardown := setup()
		defer teardown()

		type foo struct{ A string }

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			checkMethod(tc, r, "GET")
			fmt.Fprint(w, `{"A":"a"}`)
		})

		want := &foo{"a"}
		got := new(foo)

		req, _ := client.NewRequest("GET", ".", nil)
		client.Do(context.Background(), req, got)

		if !reflect.DeepEqual(got, want) {
			tc.Error("should return a response that matches the mock response", cross)
		}
	})

	t.Run("GET request that returns an HTTP error", func(tc *testing.T) {
		client, mux, _, teardown := setup()
		defer teardown()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			checkMethod(tc, r, http.MethodGet)
			w.WriteHeader(http.StatusInternalServerError)
		})

		req, _ := client.NewRequest("GET", ".", nil)
		resp, err := client.Do(context.Background(), req, nil)

		checkStatus(tc, resp, http.StatusInternalServerError)
		checkHasError(tc, err)
	})

	t.Run("GET request that receives an empty payload", func(tc *testing.T) {
		client, mux, _, teardown := setup()
		defer teardown()

		type foo struct{ A string }

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			checkMethod(tc, r, http.MethodGet)
			w.WriteHeader(http.StatusOK)
		})

		req, _ := client.NewRequest("GET", ".", nil)
		got := new(foo)
		resp, err := client.Do(context.Background(), req, got)

		checkStatus(tc, resp, http.StatusOK)
		checkNoError(tc, err)
	})

	t.Run("GET request that receives an HTML response", func(tc *testing.T) {
		client, mux, _, teardown := setup()
		defer teardown()

		type foo struct{ A string }

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			checkMethod(tc, r, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			html := `<!doctype html>
			<html lang="en-GB">
			<head>
			  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
			  <title>Default Page Title</title>
			  <link rel="shortcut icon" href="favicon.ico">
			  <link rel="icon" href="favicon.ico">
			  <link rel="stylesheet" type="text/css" href="styles.css">
			</head>
			
			<body>
			
			</body>
			</html>	`
			fmt.Fprintln(w, html)
		})

		req, _ := client.NewRequest("GET", ".", nil)
		got := new(foo)
		resp, err := client.Do(context.Background(), req, got)

		checkStatus(tc, resp, http.StatusOK)
		checkHasError(tc, err)
	})

	t.Run("request on a cancelled context", func(tc *testing.T) {
		client, mux, _, teardown := setup()
		defer teardown()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			checkMethod(tc, r, http.MethodGet)
			w.WriteHeader(http.StatusOK)
		})

		req, _ := client.NewRequest("GET", ".", nil)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		resp, err := client.Do(ctx, req, nil)

		checkHasError(tc, err)
		if resp != nil {
			tc.Error("should not return a response", cross)
		}
	})
}

// Setup establishes a test Server that can be used to provide mock responses during testing.
// It returns a pointer to a client, a mux, the server URL and a teardown function that
// must be called when testing is complete.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	c := NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	c.baseURL = url

	return c, mux, server.URL, server.Close
}

func checkMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()

	if got := r.Method; got != want {
		t.Errorf("should send a %s request to the API %s %s", want, cross, got)
	}
}

func checkHasError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Error("should return an error", cross)
	}
}

func checkNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Error("should return without error", cross, err)
	}
}

func checkStatus(t *testing.T, r *http.Response, status int) {
	t.Helper()

	if r.StatusCode != status {
		t.Errorf("should return status HTTP %d %s HTTP %d", status, cross, r.StatusCode)
	}
}
