package starling

import (
	"strings"
	"time"
)

// Errors contains a list of errors
type Errors []string

func (e Errors) Error() string {
	return strings.Join(e, ",")
}

// ErrorDetail holds the details of an error message
type ErrorDetail struct {
	Message string `json:"message"`
}

// Amount represents the value and currency of a monetary amount
type Amount struct {
	Currency   string `json:"currency"`   // ISO-4217 3 character currency code
	MinorUnits int64  `json:"minorUnits"` // Amount in the minor units of the given currency; eg pence in GBP, cents in EUR
}

// RecurrenceRule defines the pattern for recurring events
type RecurrenceRule struct {
	StartDate string `json:"startDate"`
	Frequency string `json:"frequency"`
	Interval  int32  `json:"interval,omitempty"`
	Count     int32  `json:"count,omitempty"`
	UntilDate string `json:"untilDate,omitempty"`
	WeekStart string `json:"weekStart"`
}

// Photo is a photo associated to a savings goal
type Photo struct {
	Base64EncodedPhoto string `json:"base64EncodedPhoto"` // A text (base 64) encoded picture to associate with the savings goal
}

// SpendingCategory is the category associated with a transaction
type SpendingCategory struct {
	SpendingCategory string `json:"spendingCategory"`
}

// DateRange holds two dates that represent a range. It is typically
// used when providing a range when querying the API.
type DateRange struct {
	From time.Time
	To   time.Time
}
