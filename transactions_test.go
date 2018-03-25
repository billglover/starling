package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetTransactions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mock := `{
		"transactions": [
		  {
			"id": "6f03a23a-bbfc-4479-8d4d-abb6a9119d27",
			"currency": "GBP",
			"amount": -23.45,
			"direction": "OUTBOUND",
			"created": "2017-07-05T18:27:02.335Z",
			"narrative": "Borough Barista",
			"source": "MASTER_CARD",
			"balance": 254.12
		  }
		]
	  }`
	to := time.Now()
	from := to.Add(time.Hour * -24)

	mux.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}

		params := r.URL.Query()
		if got, want := params.Get("from"), from.Format("2006-01-02"); got != want {
			t.Errorf("query string 'from': %v, want %v", got, want)
		}

		if got, want := params.Get("to"), to.Format("2006-01-02"); got != want {
			t.Errorf("query string 'to': %v, want %v", got, want)
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.GetTransactions(context.Background(), from, to)
	if err != nil {
		t.Error("unexpected error returned:", err)
	}

	want := &Transactions{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetTransactions returned %+v, want %+v", got, want)
	}
}
