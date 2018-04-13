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
	t.Log("Given the need to test that clients can be created with a default configuration:")
	t.Log("\tWhen creating a client with a default configuration:")
	c := NewClient(nil)

	if got, want := c.baseURL.String(), defaultURL; got != want {
		t.Error("\t\tshould configure the client to use the default url", cross, got)
	} else {
		t.Log("\t\tshould configure the client to use the default url", tick)
	}
	if got, want := c.userAgent, userAgent; got != want {
		t.Error("\t\tshould configure the client to use the default user-agent", cross, got)
	} else {
		t.Log("\t\tshould configure the client to use the default user-agent", tick)
	}
}

// TestNewClientWithOptions confirms that a new client can be created by passing
// in custom ClientOptions.
func TestNewClientWithOptions(t *testing.T) {
	t.Log("Given the need to test that clients can be created with options:")
	t.Log("\tWhen creating a client with a custom base URL:")

	baseURL, _ := url.Parse("https://dummyurl:4000")
	opts := ClientOptions{
		BaseURL: baseURL,
	}
	c := NewClientWithOptions(nil, opts)

	if got, want := c.baseURL.String(), baseURL.String(); got != want {
		t.Error("\t\tshould configure the client to use the custom url", cross, got)
	} else {
		t.Log("\t\tshould configure the client to use the custom url", tick)
	}
}

// TestNewRequest confirms that NewRequest returns an API request with the
// correct URL, a correctly encoded body and the correct User-Agent and
// Content-Type headers set.
func TestNewRequest(t *testing.T) {
	t.Log("Given the need to test that we can create a request:")

	c := NewClient(nil)

	t.Run("valid request", func(tc *testing.T) {

		tc.Log("\tWhen creating a:", tc.Name())
		inURL, outURL := "/foo", defaultURL+"foo"
		inBody, outBody := &topUpRequest{Amount: Amount{Currency: "GBP", MinorUnits: 1973}}, `{"amount":{"currency":"GBP","minorUnits":1973}}`+"\n"

		req, err := c.NewRequest("GET", inURL, inBody)
		checkNoError(tc, err)

		if got, want := req.URL.String(), outURL; got != want {
			t.Error("\t\tshould expand relative URLs to absolute URLs", cross, got)
		} else {
			tc.Log("\t\tshould expand relative URLs to absolute URLs", tick)
		}

		body, _ := ioutil.ReadAll(req.Body)
		if got, want := string(body), outBody; got != want {
			tc.Error("\t\tshould encode the request body as JSON", cross, got)
		} else {
			tc.Log("\t\tshould encode the request body as JSON", tick)
		}

		if got, want := req.Header.Get("User-Agent"), c.userAgent; got != want {
			tc.Error("\t\tshould pass the correct user agent", cross, got)
		} else {
			tc.Log("\t\tshould pass the correct user agent", tick)
		}

		if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
			tc.Error("\t\tshould set the content-type as application/json", cross, got)
		} else {
			tc.Log("\t\tshould set the content-type as application/json", tick)
		}
	})

	t.Run("request with invalid JSON", func(tc *testing.T) {

		tc.Log("\tWhen creating a:", tc.Name())

		type T struct{ A map[interface{}]interface{} }
		_, err := c.NewRequest("GET", ".", &T{})
		checkHasError(tc, err)
	})

	t.Run("request with an invalid URL", func(tc *testing.T) {

		tc.Log("\tWhen creating a:", tc.Name())

		_, err := c.NewRequest("GET", ":", nil)
		checkHasError(tc, err)
	})

	t.Run("request with an invalid base path", func(tc *testing.T) {

		tc.Log("\tWhen creating a:", tc.Name())

		u, _ := url.Parse("http://test.local")
		o := ClientOptions{BaseURL: u}
		lc := NewClientWithOptions(nil, o)

		_, err := lc.NewRequest("GET", "/", nil)
		checkHasError(tc, err)
	})

	t.Run("request with an invalid Method", func(tc *testing.T) {

		tc.Log("\tWhen creating a:", tc.Name())

		_, err := c.NewRequest("\n", "/", nil)
		checkHasError(tc, err)
	})

	t.Run("request with an empty body", func(tc *testing.T) {

		tc.Log("\tWhen creating a:", tc.Name())

		req, err := c.NewRequest("GET", ".", nil)
		checkNoError(tc, err)

		if req.Body != nil {
			tc.Error("\t\tshould return an empty body", cross)
		} else {
			tc.Log("\t\tshould return an empty body", tick)
		}
	})

}

// TestDo confirms that Do returns a JSON decoded value when making a request. It
// confirms the correct verb was used and that the decoded response value matches
// the expected result.
func TestDo(t *testing.T) {
	t.Log("Given the need to test that we can execute a request:")

	t.Run("successful GET request", func(tc *testing.T) {
		tc.Log("\tWhen executing a:", tc.Name())
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
			tc.Error("\t\tshould return a response that matches the mock response", cross)
		} else {
			tc.Log("\t\tshould return a response that matches the mock response", tick)
		}
	})

	t.Run("GET request that returns an HTTP error", func(tc *testing.T) {
		tc.Log("\tWhen executing a:", tc.Name())

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
		tc.Log("\tWhen executing a:", tc.Name())

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

	t.Run("request on a cancelled context", func(tc *testing.T) {
		tc.Log("\tWhen executing a:", tc.Name())

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
		} else {
			tc.Log("should not return a response", tick)
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
		t.Errorf("\t\tshould send a %s request to the API %s %s", want, cross, got)
	} else {
		t.Logf("\t\tshould send a %s request to the API %s", want, tick)
	}
}

func checkHasError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Error("\t\tshould return an error", cross)
	} else {
		t.Log("\t\tshould return an error", tick)
	}
}

func checkNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Error("\t\tshould return without error", cross, err)
	} else {
		t.Log("\t\tshould return without error", tick)
	}
}

func checkStatus(t *testing.T, r *http.Response, status int) {
	t.Helper()

	if r.StatusCode != status {
		t.Errorf("\t\tshould return status HTTP %d %s HTTP %d", status, cross, r.StatusCode)
	} else {
		t.Logf("\t\tshould return status HTTP %d %s", status, tick)
	}
}
