package uphold

import (
	"encoding/json"
	"time"

	uuid "github.com/kthomas/go.uuid"
)

// OAuthResponse is the API response returned when an authorization code has been successfully upgraded to an access token
type OAuthResponse struct {
	AccessToken  *string `json:"access_token"`
	ExpiresIn    *string `json:"expires_in"`
	TokenType    *string `json:"token_type"`
	RefreshToken *string `json:"refresh_token"`
	Scope        *string `json:"scope"`
}

// Denomination describes the value being transacted, in terms of a specific currency
type Denomination struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Pair     *string `json:"pair"`
	Rate     *string `json:"rate"`
}

// Fee describes an applied transaction fee
type Fee struct {
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Percentage *string `json:"percentage"`
	Target     *string `json:"target"`
	Type       *string `json:"type"`
}

// Destination contains properites regarding how the transaction affects the destination of the funds
type Destination struct {
	CardID      string  `json:"CardId"`      // the ID of the card credited. Only visible to the user who receives the transaction.
	Amount      float64 `json:"amount"`      // the amount credited, including commissions and fees.
	Base        float64 `json:"base"`        // the amount to credit, before commissions or fees.
	Comission   float64 `json:"comission"`   // the commission charged by Uphold to process the transaction. Commissions are only charged when currency is converted into a different denomination.
	Currency    string  `json:"currency"`    // the denomination of the funds at the time they were sent/received.
	Description *string `json:"description"` // the name of the recipient. In the case where money is sent via email, the description will contain the email address of the recipient.
	Fee         float64 `json:"fee"`         // the Bitcoin network Fee, if destination is a BTC address but origin is not.
	IsMember    *bool   `json:"isMember"`    // a boolean signaling if the destination user has completed the membership process.
	Node        *string `json:"node"`        // the details about the transaction destination node.
	Rate        *string `json:"rate"`        // the rate for conversion between origin and destination, as expressed in the currency at destination (the inverse of origin.rate).
	Type        *string `json:"type"`        //	the type of endpoint. Possible values are 'email’, 'card’ and 'external’.
}

// Origin contains properties regarding how the transaction affects the origin of the funds
type Origin struct {
	CardID      string              `json:"CardId"`      // the ID of the card debited. Only visible to the user who sends the transaction.
	Amount      float64             `json:"amount"`      // the amount debited, including commissions and fees.
	Base        float64             `json:"base"`        // the amount to debit, before commissions or fees.
	Comission   float64             `json:"comission"`   // the commission charged by Uphold to process the transaction.
	Currency    string              `json:"currency"`    // the currency of the funds at the origin.
	Description *string             `json:"description"` // the name of the sender.
	Fee         float64             `json:"fee"`         // the Bitcoin network Fee, if origin is in BTC but destination is not, or is a non-Uphold Bitcoin Address.
	IsMember    *bool               `json:"isMember"`    // a boolean signaling if the origin user has completed the membership process.
	Node        *string             `json:"node"`        // the details about the transaction origin node.
	Rate        *string             `json:"rate"`        // the rate for conversion between origin and destination, as expressed in the currency at origin (the inverse of destination.rate).
	Type        *string             `json:"type"`        //	the type of endpoint. Possible values are 'card’ and 'external’.
	Sources     []map[string]string `json:"sources"`     // the transactions where the value was originated from (id and amount).
	Username    string              `json:"username"`    // the username from the user that performed the transaction.
}

// Normalized tx property contains the normalized amount and commission values in USD
type Normalized struct {
	Amount    float64 `json:"amount"`    // the amount to be transacted.
	Comission float64 `json:"comission"` // the total commission taken on this transaction, either at origin or at destination.
	Currency  string  `json:"currency"`  // the currency in which the amount and commission are expressed. The value is always USD.
	Fee       float64 `json:"fee"`       // the normalized fee amount.
	Rate      *string `json:"rate"`      // the exchange rate for this pair.
	Target    *string `json:"type"`      //	can be origin or destination and determines where the fee was applied.
}

// Transaction represents an uphold card transaction
type Transaction struct {
	ID           *uuid.UUID       `json:"id"`           // a unique ID on the Uphold Network associated with the transaction.
	CreatedAt    *time.Time       `json:"createdAt"`    // the date and time the transaction was initiated.
	Application  *string          `json:"application"`  // the application that created the transaction.
	Denomination *Denomination    `json:"denomination"` // a message or note provided by the user at the time the transaction was initiated, with the intent of communicating additional information and context about the nature/purpose of the transaction.
	Destination  *Destination     `json:"destination"`  // the funds to be transferred, as originally requested.
	Origin       *Origin          `json:"origin"`       // the fees that were applied to the transaction.
	Normalized   *Normalized      `json:"normalized"`   // the transaction details in USD.
	Fees         []*Fee           `json:"fees"`         // the fees that were applied to the transaction.
	Message      *string          `json:"message"`      // other parameters of this transaction.
	Network      *string          `json:"network"`      // the network of the transaction (uphold for internal transactions).
	Priority     *string          `json:"priority"`     // the priority of the transaction. Possible values are normal and fast.
	Reference    *string          `json:"reference"`    // A reference assigned to the transaction.
	Params       *json.RawMessage `json:"params"`       // other parameters of this transaction.
	Status       *string          `json:"status"`       // the current status of the transaction. Possible values are: pending, waiting, cancelled or completed.
	Type         *string          `json:"type"`         // the nature of the transaction. Possible values are deposit, transfer and withdrawal.
}
