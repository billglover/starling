package starling

import (
	"io/ioutil"
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

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.userAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}
