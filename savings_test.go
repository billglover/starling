package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// TestSavingsGoals confirms that the client is able to query a list of savings goals.
func TestSavingsGoals(t *testing.T) {

	t.Log("Given the need to test fetching savings goals:")

	client, mux, _, teardown := setup()
	defer teardown()

	mock := `{
		"savingsGoalList": [
		  {
			"uid": "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b",
			"name": "Trip to Paris",
			"target": {
			  "currency": "GBP",
			  "minorUnits": 11223344
			},
			"totalSaved": {
			  "currency": "GBP",
			  "minorUnits": 11223344
			},
			"savedPercentage": 50
		  }
		]
	  }`

	t.Log("\tWhen making a call to GetSavingsGoals():")

	mux.HandleFunc("/api/v1/savings-goals", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		fmt.Fprint(w, mock)
	})

	got, _, err := client.SavingsGoals(context.Background())
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	want := &savingsGoals{}
	json.Unmarshal([]byte(mock), want)

	if got == nil {
		t.Fatal("\t\tshould not return 'nil'", cross)
	} else {
		t.Log("\t\tshould not return 'nil'", tick)
	}

	if !reflect.DeepEqual(got, &want.SavingsGoals) {
		t.Error("\t\tshould return a list of savings goals matching the mock response", cross)
	} else {
		t.Log("\t\tshould return a list of savings goals matching the mock response", tick)
	}

	if len(*got) == 0 {
		t.Errorf("\t\tshould return a list with at least one savings goal %s %d", cross, len(*got))
	} else {
		t.Log("\t\tshould return a list with at least one savings goal", tick)
	}

}

// TestGetSavingsGoals confirms that the client is able to query a single savings goal.
func TestGetSavingsGoal(t *testing.T) {

	t.Log("Given the need to test fetching an individual savings goal:")

	client, mux, _, teardown := setup()
	defer teardown()

	// Define our mock response and handler
	mock := `{
		"uid": "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b",
		"name": "Trip to Paris",
		"target": {
			"currency": "GBP",
			"minorUnits": 11223344
		},
		"totalSaved": {
			"currency": "GBP",
			"minorUnits": 11223344
		},
		"savedPercentage": 50
	}`

	t.Log("\tWhen making a call to GetSavingsGoal():")

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		fmt.Fprint(w, mock)
	})

	got, _, err := client.SavingsGoal(context.Background(), "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b")
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	want := &SavingsGoal{}
	json.Unmarshal([]byte(mock), want)

	if got == nil {
		t.Fatal("\t\tshould not return 'nil'", cross)
	} else {
		t.Log("\t\tshould not return 'nil'", tick)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return a savings goal that matches the mock response", cross)
	} else {
		t.Log("\t\tshould return a savings goal that matches the mock response", tick)
	}

	if len(got.UID) == 0 {
		t.Errorf("\t\tshould return a savings goal that has a UID %s %d", cross, len(got.UID))
	} else {
		t.Log("\t\tshould return a savings goal that has a UID", tick)
	}

}

// TestPutSavingsGoal confirms that the client is able to create a savings goal.
func TestPutSavingsGoal(t *testing.T) {

	t.Log("Given the need to test creating a savings goal:")

	client, mux, _, teardown := setup()
	defer teardown()

	uid := "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b"
	mockReq := SavingsGoalRequest{
		Name:     "test",
		Currency: "GBP",
		Target: Amount{
			Currency:   "GBP",
			MinorUnits: 10000,
		},
		Base64EncodedPhoto: "",
	}
	mockResp := `{
		"savingsGoalUid": "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b",
		"success": true,
		"errors": []
	  }`

	t.Log("\tWhen making a call to PutSavingsGoal():")

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")

		var sg = SavingsGoalRequest{}
		err := json.NewDecoder(r.Body).Decode(&sg)
		if err != nil {
			t.Fatal("\t\tshould send a request that the API can parse", cross, err)
		} else {
			t.Log("\t\tshould send a request that the API can parse", tick)
		}

		if !reflect.DeepEqual(mockReq, sg) {
			t.Error("\t\tshould send a savings goal that matches the mock", cross)
		} else {
			t.Log("\t\tshould send a savings goal that matches the mock", tick)
		}

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("\t\tshould send a savings goal with the correct UID", cross, reqUID)
		} else {
			t.Log("\t\tshould send a savings goal with the correct UID", tick)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResp)
	})

	resp, err := client.CreateSavingsGoal(context.Background(), uid, mockReq)
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("\t\tshould receive a %d status code %s %d", want, cross, got)
	} else {
		t.Logf("\t\tshould receive a %d status code %s", want, tick)
	}
}

// TestPutSavingsGoal_ValidateError confirms that the client is able to handle validation
// errors when creating savings goals.
func TestPutSavingsGoal_ValidateError(t *testing.T) {

	t.Log("Given the need to test handling validation errors when creating a savings goal:")

	client, mux, _, teardown := setup()
	defer teardown()

	uid := "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f"
	mockReq := SavingsGoalRequest{
		Name:     "test",
		Currency: "GBP",
		Target: Amount{
			Currency:   "GBP",
			MinorUnits: 10000,
		},
		Base64EncodedPhoto: "",
	}
	mockResp := `{
		"savingsGoalUid": "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f",
		"success": true,
		"errors": [
			{
				"message": "Something about the validation error"
			}
		]
	}`

	t.Log("\tWhen making a call to PutSavingsGoal():")

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")

		var sg = SavingsGoalRequest{}
		err := json.NewDecoder(r.Body).Decode(&sg)
		if err != nil {
			t.Fatal("\t\tshould send a request that the API can parse", cross, err)
		} else {
			t.Log("\t\tshould send a request that the API can parse", tick)
		}

		if !reflect.DeepEqual(mockReq, sg) {
			t.Error("\t\tshould send a savings goal that matches the mock", cross)
		} else {
			t.Log("\t\tshould send a savings goal that matches the mock", tick)
		}

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("\t\tshould send a savings goal with the correct UID", cross, reqUID)
		} else {
			t.Log("\t\tshould send a savings goal with the correct UID", tick)
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, mockResp)
	})

	resp, err := client.CreateSavingsGoal(context.Background(), uid, mockReq)
	if err == nil {
		t.Fatal("\t\texpected an error to be returned", cross)
	} else {
		t.Log("\t\texpected an error to be returned", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("\t\tshould receive a %d status code %s %d", want, cross, got)
	} else {
		t.Logf("\t\tshould receive a %d status code %s", want, tick)
	}
}

// TestAddMoney confirms that the client is able to make a request to add money to a savings goal and parse
// the successful response from the API.
func TestAddMoney(t *testing.T) {

	t.Log("Given the need to test adding money to a savings goal:")

	client, mux, _, teardown := setup()
	defer teardown()

	goalUID := "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f"
	txnUID := "28dff346-dd48-426f-96df-d7f33d29c379"
	mockResp := `{"transferUid":"28dff346-dd48-426f-96df-d7f33d29c379","success":true,"errors":[]}`

	mockAmount := Amount{Currency: "GBP", MinorUnits: 1050}
	mockReq := topUpRequest{Amount: mockAmount}

	t.Log("\tWhen making a call to AddMoney():")

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")

		var tu = topUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&tu)
		if err != nil {
			t.Fatal("\t\tshould send a request that the API can parse", cross, err)
		} else {
			t.Log("\t\tshould send a request that the API can parse", tick)
		}

		if !reflect.DeepEqual(mockReq, tu) {
			t.Error("\t\tshould send a top-up request that matches the mock", cross)
		} else {
			t.Log("\t\tshould send a top-up request that matches the mock", tick)
		}

		reqUID, err := uuid.Parse(path.Base(r.URL.Path))
		if err != nil {
			t.Error("\t\tshould send a top-up request with a valid UID", cross, reqUID)
		} else {
			t.Log("\t\tshould send a top-up request with a valid UID", tick)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResp)
	})

	id, resp, err := client.AddMoney(context.Background(), goalUID, mockAmount)
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	want := &savingsGoalTransferResponse{}
	json.Unmarshal([]byte(mockResp), want)

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("\t\tshould receive a %d status code %s %d", want, cross, got)
	} else {
		t.Logf("\t\tshould receive a %d status code %s", want, tick)
	}

	if got, want := id, txnUID; got != want {
		t.Fatal("\t\tshould be receive the UID assigned to the transaction", cross, got)
	} else {
		t.Log("\t\tshould be receive the UID assigned to the transaction", tick)
	}
}

var deleteSavingsGoalCases = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "sample savings goal",
		uid:  "840e4030-b94c-4e71-a1d3-1319a233dd3c",
	},
}

func TestDeleteSavingsGoal(t *testing.T) {
	for _, tc := range deleteSavingsGoalCases {
		t.Run(tc.name, func(st *testing.T) {
			testDeleteSavingsGoal(st, tc.name, tc.uid)
		})
	}
}

func testDeleteSavingsGoal(t *testing.T, name, uid string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodDelete)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("\t\tshould send a request with the correct UID", cross, reqUID)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.DeleteSavingsGoal(context.Background(), uid)
	checkNoError(t, err)

	if resp.StatusCode != http.StatusNoContent {
		t.Error("\t\tshould return an HTTP 204 status", cross, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	checkNoError(t, err)

	if len(body) != 0 {
		t.Error("\t\tshould return an empty body", cross, len(body))
	}
}

var savingsGoalPhotoCases = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "sample savings goal photo",
		uid:  "72011c4b-d42b-4709-8511-b7f01669e46f",
		mock: `{
			"_links": {
				 "self": {
					  "href": "api/v1/savings-goals/72011c4b-d42b-4709-8511-b7f01669e46f/photo",
					  "templated": false
				 }
			},
			"base64EncodedPhoto": "iVBORw0KGgoAAAANSUhEUgAAAPoAAABkCAYAAACvgC0OAAAAAXNSR0IArs4c6QAAHgpJREFUeAHtnQnYXdO5gBFEYo55iCSk0RiLmueZmpWr+tR8W1WUqg7aKlqlykWj3Kqhpj4t7VVjzVFJ0KiZCjVkFImEEERiyn3fOOf0/Off+5y199nn5P9jfc/zZu+99re+tda317z3+TPffFGiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB+Y1D8yePXsBWBh6w2KweBVe94KF5rVyt6I8+Gl+6AeXw2j4GMoyg5ORcBws2or0P4s28eVCcBQMh+lQFn0/Dq6DL8L87fJP2xJKKxCFXZh7S5ZYjmN/WLWE10tATyjLR5y8Ca/CKHgRxsGk+eef/wOOUao8gH/7cPk7+HJVcNLpN/HfpUk3Ylg2D+DzPYjxZ+hVJ+bT3NsOn0+ro1PYrQULs5TBEI5YEfX+MBDWgcGwJqwCi0EWmYGyTnsAu8M4jsB507MYmMd1l6Z8GwWU8fMBOlGlgQeogw6e1ul6jVwrq4Gd8LzT0Cm8jdfKtilsDDbwfmAlbFZ6Y2CzEidy/AfpXc/xGhr8exyjzDff7AAnfBygE1Uae8CGLo3kExRCnksjO0H3Wzai09hsgF+EnWBHsGE7FXeq3ipxir8t2KnsQh5+TmN/vFWJRbvRA93FA4U29NLI7RTwS7AnOB23wReaDvYaiTOIfeEL5OknHP9Ig7cHjRI98Jn0QNMNkIa0AJ5ztN4cDoOtwMbdFaQ/mbgE3Hm2scfpaVd4KjEPbfdA7oZOw3Edsi44eh8IG0JXFHftz4Me5Nl1++yumMmYp+iBVnogV0OnwWxBpo6HbWDlVmawINsrYGdIydbVBdmMZqIHuo0HMjV0GvgASmYDPwD6dptSfppRR/azKMPzjOoju1neY3ajB5rygOvrhkLj8Eu0Q1D8G5wA3a2Rl8vo7OMcyrJsOSAeowc+Cx5o2NBpFCvhCDe0rgB31BvGQacry5Zk7mTK1aMrZzLmLXqgSA/UbbQ0hg1IzM8nD4d55dtylyvOTnzfHiV64DPhgdQ1Oo3c99AXgq/O2il+vvp3mAiDYGsoupNxCv99yvgo63XTixI9ME97ILGh0wB2pdS/Ab89b6eMJ7FvgZtlfsPuj118ffcrKOJzWcxUxBHdz3Hvq4TEk+iBedQDnabuNHI/W70I2t3I3yXN8xhhb4Mp8B44qrs3cC0U/WXbItg8jPLamUSJHpinPdBhRKfS+775HPjcXCj1y6R5S226NHayNfsvhO8HfWvvN3m9P/F/D/c3aSdzdMpkJ7sM2NEsBX2gN/hMZsI74M9x/WGO52/hiw85dkmhPP7OwHIsDpbFcz9F9sMqf1o8Hd4Ay/I2TPPZcmy7kFc7+XJey/k1r27QfgzmbyqY57fJ51scu5xQjtXJlD8Uc5B8mnyOTctkpaGXKt6JKG6dptzicKfqOjhJHiPw31B0Q/ePLXyNsj+Ikz5ISrjoMNKyM10PNgEfkvsQhpmX8l6EDcDG7gN8DV6BUcQdxfFFGEN+J3Gcq0J+7JgGl9iAo29lBoI/Q+4F5fJwWimPs7QX4CniP8rxCcryOseWCmmZn7VgbSjndQ3O/aGV5VgYyjKLExu3S8nnifsUxyfgWfI6meNcF/Lkp+bngr/cVB4k7LccbyWPae2I38uxww5ToAh5EyOvZzQ0Bv3tzHGScO8EeD/B5tuETU0IDw16B8UNk9IsMow0VoVT4VGYCXlF3z4OV8NBsGy9fHJ/DXgFGsl59exU38PQ2mBZ7oPJkFf0w2NwHtgACxfsrgInwT2QtU4SpSL+NR79fgnslJZR7vnXkX4IjeQNFByRMwvxHJxeTUjgI8J8Jsn540ZPOA2aldEYcDfbTsPKcCg8AiHyIUpnpJWae2vCpJIhG+dQ8E8gbQLrgn+6xweRR36Rlm4R4WRoc/B38h/kyVydOFYW/ftT6J+UV8ILa+glW2dzfBF8XkWKI+cxUD2yJhUpKAw7i8LX4Z+QNEAQnFsmEvNKcGbWQQhraUPH/pHwGtST8dz8DnTcgyNgIDwHzchbRN4TXOdUhGt/KvoShMjdKDmF7SSEu1g/HexI1ofFwfVfRbi2gwlNC9WK/Isz15WFCjbN88HgbKWVYqN7AjaqLQBhTTd0bPQGR5HR4KjRKpmB4SugqY1g4lufb4WiGzgmO4h1bd9qn3PdkoZesvsVjs7oQuQ9lHYp581M2VhcK7rWakYeJvJdrA/czKgI109ycTN0CK8odDxZl8vEjUDszIbT4VfwFLwDHTZzuHYddROEpFWd8gAudqwOKOh8D+z8BvoVZC/NjHstTn13TlPIG079cLPwHPg99IcOHTnXRYrr6CPhf0l36TyGibcD8a6HPcFNt1bKGhg3nZYKZdLnR8EQCPWLm6P6c444tBuw/aeXTf07moaWtit8N5anBlh3E2dQgF49FTes3KnOIlaIvbJEaKTLw7Fx20AKnymkpD2T8H+l3GsmeHUifwPsTNolPgtnbpmm8ejb6C6Hlu+5lByhzx1cWiaUyTb6dTgf3DgMEQe6a+D+srJGFgN38JqVFciUs4MkGUagDTBNJnDjWrBAD6UpBYb7bX6lJwuMY75dR7trXJSciiF3edslj5DQPS1IbAo2nZW1W04iQZcLaXWqQ37Q24eAK2BAhxutvbAhXdeqJCiTbwN+CReA7TREfJV5FpzAwDu9HMFe2vVQ33JAE8fBxF0VxtfaIEHXSjpla3AaokyGEXAXOO0fj176awEUGglprICO0+WFGukm3Pez2B3h+YR7mYLIxzpE2DVTpP8o+2GQMyNHs6BKXtK/GP85whQtEzF4KThKWl/aJZb/u/A3mFQvUfy9PvfPhuXr6RV8byr2zm22zqbliTJZ/jPhOAitz9Yb/XAO+fJ1dUV8cPaA9hzNinZcY3dq6CXDt3M8GEaD53eCf4v9PY65BYfYGMz/GuAo4PvpPKKNbbB3CXmancdAVZxdOA+dZvn+/llwNHbktAOcBebHTti19xfAyrwEuMyolbI/a8ObvsYXn+CTWzC0HRxSx6Dv/F+Fp+EJmABWPGdXDgAOBJuBHWpPCJG1UDoIfp2mTN4c6U4F7WcRO9TpYH3V/x6ti/pXv68J1il9bqOrFuvHxfD36sCizimTaZ4Mx0No52qd+QlcxDPzvKNg9EQoSn6BocSMGQ59weVCU4INd7PddXf3/ZtwE0yFZuUZDPiQcwvxzdtVgRnxTcWPoe4GC/fdNPX1oq9W/gJjobz7PYHzjdMyzL2md921jR3fnphWrbxCgDvlvnFpVI6V0PEVmm8IQmU4iqkjGveOBb+lyCKjUf4f2BLsKDoJ4f5PQb4iPgHuBV9lluV2TpbqFIkAwpvadSd+H7gAsoj16ERI9ZMZ+3UWiw10fY9tT9hSIY1+MAR8X1ik+N/nuLzILcTvCXcFZupc9GpHi4ZpE2dT+BlcC3tBaufJvUIaupnC1reh/HpH318OGzTMcI0Ccey0boMQsQPfsMbEnEvCl4DhIUZKOrM4/hlSO8aUdHy96KutG+D3kFrHuZe7oRN3GbgUsojflNh5pjdyCuXo2zepcDnC/k2ce8DpT6vF6eBW4JSwSFkcYzYiP4l1apdHbHRJ0+skW/8kHafumYQ4I4kg7ZaLSNDlhcu0x8nHM3kyQLwX8PHRxL0KdmpgY1Hubw6PJ+htQVhoo30f3VPgd6TvebCg73r3TyWC42VRxB+ro+8SZc8M8aaiezL5u7pRHBt64lS7UcTS/Ukch4LrbSvehJJTOM0vFNpRrh+MwZ7rvFp5mQDTzDya1BpKuPajE6dlbybcCwmyg/goRBGdzSirv9brsHESGLftauTTtemwEk2ljy0/4bwSIzuC+yxpYl0YVHuTuHb2dhI9a+8lXPtMrgf3X5LqU0KU9gVRlsGkZiPfOUOq09A9Cf4QEsdGPiFEsUbnY64vgwvA3fJMPWSNrTnrGsLcfFoN7NF2K50fyNGNnQ7iw8I5DxLopkPIg+4Qv8GFnccKkLehO0LbAYaIo5rfiv+WMr0dEmEe03nO8sOKdcrlDGm5hPvLEGYnESJuEJ5hvQlRbpOOnaZLCQc029KWGdJ9Dd1vw42UyU6soeRt6A9h+cckkrcxzMkYhVyWEwvqunhX2BSWBsVR0VdlnRq6NxGnjWNgTShSnEK52zoqj1F8QrFmmzcbvKNRPbFzOxu2Is7lHB+FydgInRGg3q3Fzu0tqNfQLeAi+Ae3dHgb0p9wR8JG4nO4jLhjGim2+b4NdB34DmRp5OPQP5by3MYxWGzoY4K1P1W0V7yDhHI3ch6aI/d+sA3YyJN6bPO2M7rnk1bS1HY8930tUnRDXwib65HunaSbt8Hdgw173JWgkThtLc9inKU8QNpOj18BZ0t580D09gj57UVKzoLsuN3ncF2tWDbz/y68B6/DG5TpfY7KTLAhNpLFULA+WPecAWr389DT6wYykft/baAzN25bpvNgnQyJW5bMjVz7Os9ppg60goeID8cH1oxsReQLAwyshY7Ts1trdaksPO/Zhu8CVq4iZWOMWVkdcTILeXuMvF1HxO9liOyz2LaEPh4DL2HnHxzF/w7apUqXEPJl494Odgcrq53aMmAFXhjK4hTV8tjQp4BLFWdLI2AMeL+R1K7hexAhtIPXb681SmAu3Nd/WRr5ZPSPoyyZRvJyuaxcr4I97tLlwAZHG0CWDCaZc21mB7Ni0s2qsD6cbw+dGnpJZyhH7VQ39E+4Hg1OcQZAf8gqnyeCDyJXQy8ldj7HLWGL0nWWwyIomwfZFabCKBrI3Rxv4mG/wHGuCHmwznwZjoL1wJG8kdj4RV3rjp33V8Fyhcx6UOsgNvRVOoSkX6Qt/dJjdM07S5CtwZBrdrIAEe1lnRKEinH24YFvBJWelvMesCRsAqfDr6Byv8a4FfWpmrCkS9NyZzqxE6LCj+f+Q+DU3qnu1bAvOGPYB2xoP4KkqT/BqdKXO3kqYMUgeZvExbHwWCUw34kzLfOyA/wc7scfvmvVz1b4tgnpOcO6Bq6AnSGkkaOWKHbig6C6k05UTAi0XoXGcyCbF8SBx9+Yb5+rMERcFG6BrPI0Eb4FvnfeDvy10cPwPihO0awYicI9vwgL+UMMfjCxaaIRArm3MfwXJFY6wv2A4WJwpz6L7J2WZpZwEvTjkD+BO6xFil9DnQJ1Z0XcL+SDGexsCM/C3JB7SNQOb45w3gsMC5Hdy/HacSRDoR/MhOQ9SechAkOXLZUiO2J+AKMqIeEn66J6MdwLd8A5sBk47VScatjrp8kD3AjpbV33bU7hEkcvRk4/OrkBJiclRLhT+Rsgy6xFU/38p1khfWcvR4Gj8TPN2quKvyTnZ4Gju1P8lgn218f4RbB2yxLJbvjjwCiVDiJQv6urbU4GT+OZ2L6CZQEqohtxjwbH6Kzo2qvcuKvvGub3z72rA6vOH+H8yarrtFN3bZeHpDTS4tSGuxkzvTawwfXq5N2OsGnBx/7p6jMx9CU4FG6ErPkhSqI48/Czzi8k3m0yELuuhX8LWzRpqsjodt7vBBq07sxrcjAFOjNL/XRjRXkRpkHiWliFnOIabEMYURufiu+0fSjhbja5/qgWe2vX8HeBOo6E70NeccTz1U8WsYK4FixMKPMEjF1LuZ0BDYQ9YSdYD2p9QFCwrIPmOdg9kDSK6kBcFjmLOhKcqYXKLBQdOJ6DcTADFAeE1aA/OPVcGfJ2pDZ095ZC5HMhSt1Qx1miA+WVIXlfsKT0OsenYduQSBl0fJhW5BEpce4k/AfgqOHrF6fXt4MNwcrin4tyxpFbqKx9iPx1WCGjkaUy6gerU6apKLv34EbduWBHtA240WKjt1PqCVnEZ3c4DMkSqYGuz+W/G+iUb9vBXA9Xgx2zS0Kf3WxQbNTWN6fSS4BLP2c4R4CdQBZxIBgdGME9nB74PHSqH2i2JWr3Y9W6cWCA9d7ouOntZ+JDA/TnbGj55ZE/22uFuGmSOlpxz5/kuc70V1hZR93E8mFnQVgVdgE3wvLIUCK1dX1Hen4c0B+OgZtgPGSR4Si7dq8I17k344h7BPh79EYyGoV9K4kGnhBnZXgKGkmHzTjNE8E3PyF5s0N1j6EtQlp5N+P86av+WAkehFAZiWL48gTlr0GI40IzUNbzN8z23onCvYUTb2QMxI4O1lG7wRkwAmZBXvENQtZRNWOu09VJ2/JsA3aCUyBE7Bh2qLbKdTMN/aqARG1I+1WnGXpOPDtj3940kqSGvjaRXm8UkfsfwjmheWpWj7TyNHQ79f7ltDnfDsZBqPh3IOruYVWvkZxuTS4nVuDR6V/qGo9ph1O8IsSR7CdwM/wUtoRmOhHzVZ52ctpewS+fwDA4mpQPh9cCcuAeSyE741QcfTcoIM0r0bklQC9JRf/m9fFE4g5LMloT5nLhEMqzY014V7i07H+AI3jOY6oy9ADn3wE/ZAuRk1DybwVY1kSpbuhj0XglUStfoJsyD8MZcH8+E5li6RQ7KitoEeKeQZdY11EJyvsWjcrlDMTGXoS4bu6wDEgx6n89nddP7k30T7HbKPhtFNysDZGVUDqNhhA+xQ2x2pzOJ0S3kftHHKdVm+LaDsDO8y8Q0hE6mv8IdoNEqTR0jL+FxlMQYjjRGIE+cGcGZ4I9qAmfie2XOOYWHlBPWBF8XbdWkiHScOPHDuWdpPs5wvRHU0JeB8MPwBHFtVdqjxuQ0Hh0rBz1xOdZ1L6CtkLym6uR4wvXzWeBG3OZheetL4bDuMDIW6Hnm4m+gfqd1IjbB/widLVON7MH2FGdRTneSIpaqs/OUB0sQ8RO+TzyNjhJudLQSzdHcMxbwXX81bAbmTwV/Cst0+Gjku1MBzJs4+4HuxPxQngAboZvQJo8z42mOpUqwzas3J0e+V6V+NfBL+EauA9OI3xLWB7m5zpI0F0URZcitc+rNr6dnRWoCHFGE9Jp+sVccFnMGPqbcLgMPDYjLxJ5KIQ8J/N4OFxN+tuCu9ZBgu5y4KDlCPsg3MH11kGR05XM88z027zb5Y9zcN+3UpPq6VXdW5NzN9WXqgrrfIrCavAE5BF/adWvs9VsIdhYHfYAv5V/FD6CavGXYamjAPcug4+rI+Q8PzJbzjtqk+a5Ken6hwz9e2mnwI4wCBLLQ/hC4KaTI9EsaCRvoHBwdU64bmYz7tZGCXL/ZdihOs20c/SWhUPhBcginTbjymlgZAswD1nETTz/VqI/gx4Ai5TteeRav7ux6+8JjoGbYQZUi/WwU30nLHQzzme1enW6aefoHQshzx+1Oe3FutehTB2mZvQg7vQNI0F3yXukJZwS/gjxx6bcaxhMujrtp7AVDIC0KehA7u0Bf4Qkscf1PWTI+jIpvmHvw5NpNxuFUxZ7+0NT9GzU5l9MZyJMJo69tz23Myp7+6XBsq4Nq0GIaCt0qhdibwRKezZQtLK6O38xxz9TB16p1ifctb4jjZ3B/rAxZK1bREkW0vPb77O5OwR6JWt1Cl2OkG+Ds0PrrG+GpnJ8D3pCH1gFnJV5niQbEuir4SPJg8+slXIpxs3PKQGJ6NvjwXX/WWX9Dg29FHgHRytp/eG/bOE/R6dGzYgNW+cNamBkce7vDH9M0fuEcBtKM/IckSfkMcCD19FHgZWpkVgx1yhR1jX/ygKfHoL/NZ6d3LjgGI0V70LFTZ7EGUdVdNe9Z8BhlP/fHMeDfrAOrQY2mBXBRtQKuQ6jG8HRkKUeOuqtWYJDZvGDn8PhwswxM0SgI3FWO4Qo7k/tExBVP/vnn+8j7kj1kyrTMMJ9WFllXQzb6yQK9/wox7VpUudiHCvHnYmROwb6IP3F3PodgytXm3PWzGiuISt43l7aCrcFZKlwqFfEZ5L0XCoKKSeO5v6F03JHkaKWKXgU2n8NjGHlGgxWxGPhm3AwbAnO1rzfEqHMrnXtkG6EZjv5LHm0TF+lLoZ06lnsdtKljM72ToexnW4mB5ins8mbvu9coTA4g3A3vbKKU0zf5dlzzxHOe4Fr7oMIuARsQF+cc7PmH9KdRdBw8NhI7Nl+iN0Vyoqc94a9uN4f8jYyzbmZ5f8K+4EXOWR14lTylSN+nigfEuk88vx4nshpcUrPxNHKTjiL6P/QZ/BxFsNpuuR1GveOg+vhozS9FoQ7Y2nL86aMT5LWuWAbDZHtUTqddrFY2uh6Ewong+vELGKc9TF8B0fXO+vBNoaB4mizO/zDiwR5mrCXwYbcSL6Cgn/A/3aOpuVIegAsD83ILUQ2H3nlESLaoe0NLRvFqjJnQ7kMLq0KK/JUX/wYLoBlijSMrRFgI/lcEXZpCJOoD84m7Jhs9L2KsFvHhr53Fjqujk7Rty7HoHtoRwcaPgy9hxN1cVYPuBRaIY9gdKmkhAmfH65pRaKBNt0JdarZlGDDDmgv8Dt7/yeNVsmrGD4eXGsmCvdy77qXDWLDnWTL8xoUIe4gO60cCMMDDN6Ljns4wYL+vjA0wHZelceJeDgsUZ0prvWVs81GYl1z9pdZiLcU3NIogar7ZySuBekZ7an+AKFThCyZHYDyZkkRSNf11etJ99oUdivpOCI3JZTD7we09S3YA5xuOVMpav2onfvhIPA/JZjJMU2cQvdIu1kVnqqDfT/HtTyHguk2I88S+Uj4OfisQ2Y9LqMyTcfJr7PSr8GJ0MwMjegd5AWunOEcQBpXwfQOdz+9CFm2qBOi18k8ab5FoHsS1ql68hI3vw9DUpXoDZaE66t6haJOP8HQhUkJE+6I/teiEspox5nGwKR8NRuG3fJ72b05HwIjYQr4brb2OwGCOok+ew8cUf0BxEGwVEi+0PNjD+M4ivoDjw8ScNZho2go6DmaOJI9ANNAe/XENLXvn6E6EVYqJ8K5M8crYCb47UMSfnfwvXKcrEfiOsLqgy/D/8E4eBca5RuVOc9Gv08E36V/BVaE1E7R/HF/PyinkeRvn4X+WzJrear1S+m8ybEspuVM4T44AVaBOXmt26Og5Hr6BvBdaJFiT+MfSXiy2ijp7cS1M4lm19nVZkPOx6J0CPkZHqLcrA7ldKo9ADaEwdAfLLON13s+nFnwLrwBrgHd1/DPZum7TEJ6qxLB/Qufo/skZSk//ycIuBvbH5VvNDpic2F03HvxmVmGVaCcf2eE5t03Ac/BQzAS+9M4dhDsWO79YRDU5s2Zy1PgnwrTH00JaVnelWEDMO/uDZh+H1gUnOGadzdkJ4G+doPznzCWPJifhkI6LjOcya0DxqmOZxrvwM3YG8Mxt5CO9WRPOAqc9ZjPe+FpbH/IsSLlB10JqD7BkJtx58Ph1eEFnT+AHacUToXMsM4/oXTk0DbxgfpH8W9sW4oJCZUqR29uLQz6wwc3k3zN4Njlhfw7BTf/Hu0wzLuNvUtLqbHYyM13uaHP6C5+D3Vu3YauERyxKYfbYFmvCxZ7z4lgxbanbbc8SoLf5aEOa3fCMb3ogXZ6wB6skTyGghsx1dOPRnFC79vA+8LcaOQPku5RsZGHPqqo15090HBEt3CM6mtycGq7ltfdXFx/XQEX0MgndPOyxOxHDxTrARq7O46tfCeM+ZaKO7u+s927WM9Ea9EDXd8DQSO6xaCBqHsufNfrbibu7vv24ApG8de7Wd5jdqMH2usBGrvfk/8Guov8nYx+AyrvbtvrsZha9EA39QCNxl+gXQNdVfxQ4Tb4KvhuN0r0wGfeA8FT92pP0YB81XYRHAALVt+bS+e+a3Zj7R64Bp5jiv4uxyjRA9EDeCBXQ9dzNPYlOLhePwaWg3aLDXk8+C7cX7ANpXFP4RgleiB6oMYDuRt62Q4N3s8XT4TNodWjuw15DPwLHgY/2/Q6SvRA9EAdDzTd0LVNY/cb5+3gCNgBFoMi5E2MPA8j4QkYBaNp3H7/HSV6IHog0AOFNPRyWjR4p/D+wGFb2BrWBsP8yL+e+G30O+CrLxv2c/AMvAiTYQqNu95PMVGJEj0QPZDmgUIbejkRGrxTeH+F5a+C+sNAGARrgD/aMN2p8AqMhdfAzTTD3ocPaNgdfn1DWJTogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogeiB6IHogbnsgf8HTLgkC6jdWKgAAAAASUVORK5CYII="
	  }`,
	},
}

func TestSavingsGoalPhoto(t *testing.T) {
	for _, tc := range savingsGoalPhotoCases {
		t.Run(tc.name, func(st *testing.T) {
			testSavingsGoalPhoto(st, tc.name, tc.uid, tc.mock)
		})
	}
}

func testSavingsGoalPhoto(t *testing.T, name, uid string, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		resource := path.Base(r.URL.Path)
		reqUID := path.Base(path.Dir(r.URL.Path))

		if reqUID != uid {
			t.Error("\t\tshould send a request with the correct UID", cross, reqUID)
		}

		if resource != "photo" {
			t.Error("\t\tshould request the photo", cross, resource)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mock)
	})

	photo, resp, err := client.SavingsGoalPhoto(context.Background(), uid)
	checkNoError(t, err)

	if resp.StatusCode != http.StatusOK {
		t.Error("\t\tshould return an HTTP 200 status", cross, resp.Status)
	}

	want := &Photo{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(photo, want) {
		t.Error("\t\tshould return a savings goal photo that matches the mock", cross)
	}

	if len(photo.Base64EncodedPhoto) == 0 {
		t.Error("\t\tshould return a base64 encoded photo", cross)
	}
}

// TestWithdraw confirms that the client is able to make a request to withdraw money from a savings goal.
func TestWithdraw(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	goalUID := "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f"
	txnUID := "28dff346-dd48-426f-96df-d7f33d29c379"
	mockResp := `{"transferUid":"28dff346-dd48-426f-96df-d7f33d29c379","success":true,"errors":[]}`

	mockAmount := Amount{Currency: "GBP", MinorUnits: 1050}
	mockReq := withdrawalRequest{Amount: mockAmount}

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")

		var wr = withdrawalRequest{}
		err := json.NewDecoder(r.Body).Decode(&wr)
		if err != nil {
			t.Fatal("should send a request that the API can parse", cross, err)
		}

		if !reflect.DeepEqual(mockReq, wr) {
			t.Error("should send a top-up request that matches the mock", cross)
		}

		resource := path.Base(path.Dir(r.URL.Path))
		if resource != "withdraw-money" {
			t.Error("should make a request to withdraw-money", cross, resource)
		}

		reqUID, err := uuid.Parse(path.Base(r.URL.Path))
		if err != nil {
			t.Error("should send a top-up request with a valid UID", cross, reqUID)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResp)
	})

	id, resp, err := client.Withdraw(context.Background(), goalUID, mockAmount)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	want := &savingsGoalTransferResponse{}
	json.Unmarshal([]byte(mockResp), want)

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("should receive a %d status code %s %d", want, cross, got)
	}

	if got, want := id, txnUID; got != want {
		t.Fatal("\t\tshould be receive the UID assigned to the transaction", cross, got)
	}
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	goalUID := "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f"
	mockResp := `["INSUFFICIENT_FUNDS"]`

	mockAmount := Amount{Currency: "GBP", MinorUnits: 10000000}
	mockReq := withdrawalRequest{Amount: mockAmount}

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")

		var wr = withdrawalRequest{}
		err := json.NewDecoder(r.Body).Decode(&wr)
		if err != nil {
			t.Fatal("should send a request that the API can parse", cross, err)
		}

		if !reflect.DeepEqual(mockReq, wr) {
			t.Error("should send a top-up request that matches the mock", cross)
		}

		resource := path.Base(path.Dir(r.URL.Path))
		if resource != "withdraw-money" {
			t.Error("should make a request to withdraw-money", cross, resource)
		}

		reqUID, err := uuid.Parse(path.Base(r.URL.Path))
		if err != nil {
			t.Error("should send a top-up request with a valid UID", cross, reqUID)
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, mockResp)
	})

	_, resp, err := client.Withdraw(context.Background(), goalUID, mockAmount)
	if err == nil {
		t.Fatal("should return an error when making the request", cross)
	}

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("should receive a %d status code %s %d", want, cross, got)
	}
}
