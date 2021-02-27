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

var feedTC = []struct {
	name  string
	act   string
	cat   string
	since time.Time
	mock  string
}{
	{
		name: "no transactions",
		act:  "30aa7ab8-4389-4658-a4f8-0bc6d0015ba0",
		cat:  "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
		since: time.Now(),
		mock: `{
		"feedItems": []
		}`,
	},
	{
		name: "single transaction",
		act:  "30aa7ab8-4389-4658-a4f8-0bc6d0015ba0",
		cat:  "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
		since: time.Now(),
		mock: `{
		"feedItems": [
			 {
				  "feedItemUid": "dbb59f1c-39e6-4558-87ba-11c142965393",
				  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
				  "amount": {
						"currency": "GBP",
						"minorUnits": 32
				  },
				  "sourceAmount": {
						"currency": "GBP",
						"minorUnits": 32
				  },
				  "direction": "OUT",
				  "transactionTime": "2018-06-28T07:16:28.364Z",
				  "source": "MASTER_CARD",
				  "sourceSubType": "CHIP_AND_PIN",
				  "status": "SETTLED",
				  "counterPartyType": "MERCHANT",
				  "counterPartyUid": "e6dbe57e-7c23-4015-97a4-4afbbf7faa23",
				  "reference": "ATM 111072\\35 REGENT ST), LONDON\\LONDON\\SW1Y 4ND  00 GBR",
				  "country": "GB",
				  "spendingCategory": "HOLIDAYS"
			 }
			 ]
		}`,
	},
	{
		name: "multiple transactions",
		act:  "30aa7ab8-4389-4658-a4f8-0bc6d0015ba0",
		cat:  "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
		since: time.Now(),
		mock: `{
			"feedItems": [
				 {
					  "feedItemUid": "dbb59f1c-39e6-4558-87ba-11c142965393",
					  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
					  "amount": {
							"currency": "GBP",
							"minorUnits": 32
					  },
					  "sourceAmount": {
							"currency": "GBP",
							"minorUnits": 32
					  },
					  "direction": "OUT",
					  "transactionTime": "2018-06-28T07:16:28.364Z",
					  "source": "MASTER_CARD",
					  "sourceSubType": "CHIP_AND_PIN",
					  "status": "SETTLED",
					  "counterPartyType": "MERCHANT",
					  "counterPartyUid": "e6dbe57e-7c23-4015-97a4-4afbbf7faa23",
					  "reference": "ATM 111072\\35 REGENT ST), LONDON\\LONDON\\SW1Y 4ND  00 GBR",
					  "country": "GB",
					  "spendingCategory": "HOLIDAYS"
				 },
				 {
					  "feedItemUid": "199c2bba-9f4d-4b20-b5df-4de440411b03",
					  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
					  "amount": {
							"currency": "GBP",
							"minorUnits": 7
					  },
					  "sourceAmount": {
							"currency": "GBP",
							"minorUnits": 7
					  },
					  "direction": "OUT",
					  "transactionTime": "2018-06-28T07:16:28.361Z",
					  "source": "MASTER_CARD",
					  "sourceSubType": "CHIP_AND_PIN",
					  "status": "SETTLED",
					  "counterPartyType": "MERCHANT",
					  "counterPartyUid": "c052f76f-e919-427d-85fc-f46a75a3ff26",
					  "reference": "MASTERCARD EUROPE      WATERLOO      BEL",
					  "country": "GB",
					  "spendingCategory": "HOLIDAYS"
				 },
				 {
					  "feedItemUid": "32f8ffc4-d12c-43fe-9d1b-61faf7243143",
					  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
					  "amount": {
							"currency": "GBP",
							"minorUnits": 24
					  },
					  "sourceAmount": {
							"currency": "GBP",
							"minorUnits": 24
					  },
					  "direction": "OUT",
					  "transactionTime": "2018-06-28T07:16:28.359Z",
					  "source": "MASTER_CARD",
					  "sourceSubType": "CHIP_AND_PIN",
					  "status": "SETTLED",
					  "counterPartyType": "MERCHANT",
					  "counterPartyUid": "c052f76f-e919-427d-85fc-f46a75a3ff26",
					  "reference": "MASTERCARD UK MANA\\19TH FLOOR\\LONDON E14\\E14 5NP      GBR",
					  "country": "GB",
					  "spendingCategory": "HOLIDAYS"
				 }
			]
	  }`,
	},
	{
		name:  "transaction since",
		act:   "30aa7ab8-4389-4658-a4f8-0bc6d0015ba0",
		cat:   "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
		since: time.Now(),
		mock: `{
		"feedItems": [
			 {
				  "feedItemUid": "dbb59f1c-39e6-4558-87ba-11c142965393",
				  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
				  "amount": {
						"currency": "GBP",
						"minorUnits": 32
				  },
				  "sourceAmount": {
						"currency": "GBP",
						"minorUnits": 32
				  },
				  "direction": "OUT",
				  "transactionTime": "2018-06-28T07:16:28.364Z",
				  "source": "MASTER_CARD",
				  "sourceSubType": "CHIP_AND_PIN",
				  "status": "SETTLED",
				  "counterPartyType": "MERCHANT",
				  "counterPartyUid": "e6dbe57e-7c23-4015-97a4-4afbbf7faa23",
				  "reference": "ATM 111072\\35 REGENT ST), LONDON\\LONDON\\SW1Y 4ND  00 GBR",
				  "country": "GB",
				  "spendingCategory": "HOLIDAYS"
			 }
			 ]
		}`,
	},
}

func TestFeed(t *testing.T) {
	for _, tc := range feedTC {
		t.Run(tc.name, func(t *testing.T) {
			testFeed(t, tc.name, tc.act, tc.cat, tc.mock, tc.since)
		})
	}
}

func testFeed(t *testing.T, name, act, cat, mock string, since time.Time) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/feed/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		if r.URL.Path != "/api/v2/feed/account/"+act+"/category/"+cat {
			t.Error("should sent a request to the correct path")
		}

		params := r.URL.Query()

		if time.Time.IsZero(since) == false {
			if got, want := params.Get("changesSince"), since.Format(time.RFC3339Nano); got != want {
				t.Errorf("should include 'changesSince=%s' query string parameter %s 'changesSince=%s'", want, cross, got)
			}
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.Feed(context.Background(), act, cat, since)
	checkNoError(t, err)

	want := &feed{}
	json.Unmarshal([]byte(mock), want)

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if len(got) != len(want.Items) {
		t.Error("should return the correct number of items")
	}
	if !reflect.DeepEqual(got, want.Items) {
		t.Error("should return a slice of feed items matching the mock response")
	}
}

func TestFeedForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/feed/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.Feed(context.Background(), "30aa7ab8-4389-4658-a4f8-0bc6d0015ba0", "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0", time.Now())
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a slice of items")
	}
}

var itemTC = []struct {
	name string
	act  string
	cat  string
	itm  string
	mock string
}{
	{
		name: "card transaction",
		act:  "30aa7ab8-4389-4658-a4f8-0bc6d0015ba0",
		cat:  "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
		itm:  "dbb59f1c-39e6-4558-87ba-11c142965393",
		mock: `{
			"feedItemUid": "dbb59f1c-39e6-4558-87ba-11c142965393",
			"categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
			"amount": {
				"currency": "GBP",
				"minorUnits": 32
			},
			"sourceAmount": {
				"currency": "GBP",
				"minorUnits": 32
			},
			"direction": "OUT",
			"transactionTime": "2018-06-28T07:16:28.364Z",
			"source": "MASTER_CARD",
			"sourceSubType": "CHIP_AND_PIN",
			"status": "SETTLED",
			"counterPartyType": "MERCHANT",
			"counterPartyUid": "e6dbe57e-7c23-4015-97a4-4afbbf7faa23",
			"reference": "ATM 111072\\35 REGENT ST), LONDON\\LONDON\\SW1Y 4ND  00 GBR",
			"country": "GB",
			"spendingCategory": "HOLIDAYS"
		}`,
	},
}

func TestFeedItem(t *testing.T) {
	for _, tc := range itemTC {
		t.Run(tc.name, func(t *testing.T) {
			testFeedItem(t, tc.name, tc.act, tc.cat, tc.itm, tc.mock)
		})
	}
}

func testFeedItem(t *testing.T, name, act, cat, itm, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/feed/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		if r.URL.Path != "/api/v2/feed/account/"+act+"/category/"+cat+"/"+itm {
			t.Error("should sent a request to the correct path")
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.FeedItem(context.Background(), act, cat, itm)
	checkNoError(t, err)

	want := new(Item)
	json.Unmarshal([]byte(mock), want)

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("should return an item matching the mock response")
	}

	if got.FeedItemUID != want.FeedItemUID {
		t.Error("should return the requested item")
	}
}

func TestFeedItemForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/feed/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.FeedItem(context.Background(), "30aa7ab8-4389-4658-a4f8-0bc6d0015ba0", "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0", "dbb59f1c-39e6-4558-87ba-11c142965393")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return an item")
	}
}
