package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)

	n := Notification{}
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		color.Red(err.Error())
	}

	if n.Direction == "PAYMENT" {
		printTxn(n)
	}
}

type Notification struct {
	WebhookType string  `json:"webhookType"`
	Direction   string  `json:"direction"`
	Content     Content `json:"content"`
}

type Content struct {
	Amount         float32 `json:"amount"`
	SourceCurrency string  `json:"sourceCurrency"`
	ForCustomer    string  `json:"forCustomer"`
	CounterParty   string  `json:"counterParty"`
}

func printTxn(n Notification) {
	ru := roundUp(n.Content.Amount)
	fmt.Printf("You just spent ")
	color.Set(color.FgYellow)
	fmt.Printf("%.2f%s", -n.Content.Amount, n.Content.SourceCurrency)
	color.Unset()
	fmt.Printf(" at %s, you could have donated ", n.Content.CounterParty)
	color.Set(color.FgGreen)
	fmt.Printf("%.2f%s\n", ru, n.Content.SourceCurrency)
	color.Unset()
}

func roundUp(a float32) float32 {
	if a < 0 {
		a = -a
	}

	f := int(a*100) - (int(a) * 100)
	r := float32(100-f) / 100
	return r
}
