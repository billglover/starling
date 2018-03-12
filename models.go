package starling

type ErrorDetail struct {
	Message string // The error message
}

// SavingsGoalTransferResponse represents the response received after attempting to make an immediate or recurring transfer
// into/out of a savings goal.
type SavingsGoalTransferResponse struct {
	UID     string        `json:"transferUid"` // Unique identifier for the transfer
	Success bool          `json:"success"`     // True if the method completed successfully
	Errors  []ErrorDetail `json:"errors"`      // List of errors if the method request failed
}

type CurrencyAndAmount struct {
	Currency   string `json:"currency"`   // ISO-4217 3 character currency code
	MinorUnits int64  `json:"minorUnits"` // Amount in the minor units of the given currency; eg pence in GBP, cents in EUR
}

// TopUpRequest represents request to make an immediate transfer into a savings goal
type TopUpRequest struct {
	Amount CurrencyAndAmount `json:"amount"`
}

type RecurrenceRule struct {
	StartDate string `json:"startDate"`
	Frequency string `json:"frequency"`
	Interval  int32  `json:"interval,omitempty"`
	Count     int32  `json:"count,omitempty"`
	UntilDate string `json:"untilDate,omitempty"`
	WeekStart string `json:"weekStart"`
}

// ScheduledSavingsPaymentRequest represents a request to create scheduled payment into a savings goal
type ScheduledSavingsPaymentRequest struct {
	RecurrenceRule    RecurrenceRule    `json:"recurrenceRule"`
	CurrencyAndAmount CurrencyAndAmount `json:"currencyAndAmount"`
}

// CreateOrUpdateSavingsGoalResponse represents a response after attempting to create a savings goal
type CreateOrUpdateSavingsGoalResponse struct {
	SavingsGoalUid string        `json:"savingsGoalUid"`
	success        bool          `json:"success"`
	errors         []ErrorDetail `json:"errors"`
}

// SavingsGoalRequest is a request to create a new savings goal
type SavingsGoalRequest struct {
	Name               string            `json:"name"`     // Name of the savings goal
	Currency           string            `json:"currency"` // ISO-4217 3 character currency code of the savings goal
	Target             CurrencyAndAmount `json:"target"`
	Base64EncodedPhoto string            `json:"base64EncodedPhoto"` // A text (base 64) encoded picture to associate with the savings goal
}

// SavingsGoal is a goal defined by a customer to hold savings
type SavingsGoal struct {
	UID             string            `json:"uid"`  // Unique identifier of the savings goal
	Name            string            `json:"name"` // Name of the savings goal
	Target          CurrencyAndAmount `json:"target"`
	TotalSaved      CurrencyAndAmount `json:"totalSaved"`
	SavedPercentage int32             `json:"savedPercentage"` // Percentage of target currently deposited in the savings goal
}

// SavingsGoals is a list containing all savings goals for customer
type SavingsGoals struct {
	SavingsGoalList []SavingsGoal `json:"savingsGoalList"`
}

// SavingsGoalPhoto is a photo associated to a savings goal
type SavingsGoalPhoto struct {
	base64EncodedPhoto string `json:"base64EncodedPhoto"` // A text (base 64) encoded picture to associate with the savings goal
}

// WithdrawalRequest is a request to withdraw money from a savings goal
type WithdrawalRequest struct {
	Amount CurrencyAndAmount `json:"amount"`
}

// MastercardTransactionPayload is the webhook payload for mastercard transactions
type MastercardTransactionPayload struct {
	WebhookNotificationUID string            `json:"webhookNotificationUid"` // Unique identifier of the webhook dispatch event
	CustomerUID            string            `json:"customerUid"`            // Unique identifier of the customer
	WebhookType            string            `json:"webhookType"`            // The type of the event
	EventUID               string            `json:"eventUid"`               // Unique identifier of the customer transaction event
	TransactionAmount      CurrencyAndAmount `json:"transactionAmount"`
	SourceAmount           CurrencyAndAmount `json:"sourceAmount"`
	Direction              string            `json:"direction"`            // The cashflow direction of the card transaction
	Description            string            `json:"description"`          // The transaction description, usually the name of the merchant
	MerchantUID            string            `json:"merchantUid"`          // The unique identifier of the merchant
	MerchantLocationUID    string            `json:"merchantLocationUid"`  // The unique identifier of the merchant location
	Status                 string            `json:"status"`               // The status of the transaction
	TransactionMethod      string            `json:"transactionMethod"`    // The method of card usage
	TransactionTimestamp   string            `json:"transactionTimestamp"` // Timestamp of the card transaction
	MerchantPosData        MerchantPosData   `json:"merchantPosData"`
}

// MerchantPosData is data relating to the merchant at the point-of-sale
type MerchantPosData struct {
	PosTimestamp       string `json:"posTimestamp"`       // The transaction time as reported at the point of sale
	CardLast4          string `json:"cardLast4"`          // The last 4 digits of the mastercard PAN
	AuthorisationCode  string `json:"authorisationCode"`  // The authorisation code of the transaction, as reported at the point of sale
	Country            string `json:"country"`            // The country of the transaction, in ISO-3 format
	MerchantIdentifier string `json:"merchantIdentifier"` // The merchant identifier as reported by Mastercard AKA mid
}

// Account represents bank account details
type Account struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	AccountNUmber string `json:"accountNumber"`
	SortCode      string `json:"sortCode"`
	Currency      string `json:"currency"`
	IBAN          string `json:"iban"`
	BIC           string `json:"bic"`
	CreatedAt     string `json:"createdAt"`
}
