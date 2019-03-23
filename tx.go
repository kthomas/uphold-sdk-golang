package uphold

import "fmt"

// CommitTransaction commits a previously quoted transaction
func CommitTransaction(token, cardID, transactionID string) (*Transaction, error) {
	var tx *Transaction
	var err error

	client, err := NewUpholdAPIClient(stringOrNil(token), stringOrNil("/v0/me/"))
	if err != nil {
		return nil, err
	}

	status, err := client.Post(fmt.Sprintf("cards/%s/transactions/%s/commit", cardID, transactionID), nil, &tx)
	if err != nil {
		log.Warningf("Failed to authorize client credentials on behalf of client id: %s; %s", upholdClientID, err.Error())
		return nil, err
	}

	log.Debugf("Received %d status code when attempting to commit transaction (tx id: %s) on behalf of client id: %s; response: %s", status, transactionID, upholdClientID, tx)

	return tx, err
}

// CreateTransaction submits a transaction to the Uphold platform but does not commit it for settlement
func CreateTransaction(token, cardID, currency, destination string, amount float64) (*Transaction, error) {
	var tx *Transaction
	var err error

	client, err := NewUpholdAPIClient(stringOrNil(token), stringOrNil("/v0/me/"))
	if err != nil {
		return nil, err
	}

	status, err := client.Post(fmt.Sprintf("cards/%s/transactions", cardID), map[string]interface{}{
		"demonination": map[string]interface{}{
			"amount":   amount,
			"currency": currency,
		},
		"destination": destination,
	}, &tx)
	if err != nil {
		log.Warningf("Failed to authorize client credentials on behalf of client id: %s; %s", upholdClientID, err.Error())
		return nil, err
	}

	log.Debugf("Received %d status code in response to attempted transaction creation API call on behalf of client id: %s; response: %s", status, upholdClientID, tx)

	return tx, err
}
