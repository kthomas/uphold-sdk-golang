package uphold

import "fmt"

// CreateUser creates a new Uphold user
func CreateUser(email, password string, country, locale, accountType *string) (*User, error) {
	var user *User
	var err error

	client, err := NewUnauthorizedAPIClient(stringOrNil("/v0/users"))
	if err != nil {
		return nil, err
	}

	if country == nil {
		country = stringOrNil("US")
	}

	if locale == nil {
		locale = stringOrNil("en-US")
	}

	if accountType == nil {
		accountType = stringOrNil("business")
	}

	status, err := client.Post("", map[string]interface{}{
		"country":  country,
		"email":    email,
		"password": password,
		"type":     accountType,
		"settings": map[string]interface{}{
			"hasMarketingConsent": false,
		},
		"intl": map[string]interface{}{
			"dateTimeFormat": map[string]interface{}{
				"locale": locale,
			},
			"language": map[string]interface{}{
				"locale": locale,
			},
			"numberFormat": map[string]interface{}{
				"locale": locale,
			},
		},
	}, &user)
	if err != nil {
		log.Warningf("Failed to create uphold user; %s", upholdClientID, err.Error())
		return nil, err
	}

	if status == 200 {
		log.Debugf("Received %d status code when attempting to create uphold user; response: %s", status, user)
		return user, nil
	}

	return nil, fmt.Errorf("Failed to create uphold user; status code: %d", status)
}

// CreateDocument upserts a document on behalf of an uphold account holder
func CreateDocument(token, documentType string, value interface{}) error {
	var resp map[string]interface{}
	var err error

	client, err := NewUpholdAPIClient(stringOrNil(token), stringOrNil("/v0/me"))
	if err != nil {
		return err
	}

	status, err := client.Post("documents", map[string]interface{}{
		"type":  documentType,
		"value": value,
	}, &resp)
	if err != nil {
		log.Warningf("Failed to create document on behalf of uphold user; %s", err.Error())
		return err
	}

	if status == 200 {
		log.Warningf("Created document on behalf of uphold user; %s", err.Error())
		return nil
	}

	return fmt.Errorf("Failed to create document on behalf of uphold user; status: %d", status)
}

// AddPhone adds a phone to an uphold account
func AddPhone(token, countryCode, phone string) error {
	var resp map[string]interface{}
	var err error

	client, err := NewUpholdAPIClient(stringOrNil(token), stringOrNil("/v0/me"))
	if err != nil {
		return err
	}

	status, err := client.Post("phones", map[string]interface{}{
		"countryCode": countryCode,
		"phone":       phone,
	}, &resp)
	if err != nil {
		log.Warningf("Failed to create document on behalf of uphold user; %s", err.Error())
		return err
	}

	if status == 200 {
		log.Warningf("Created document on behalf of uphold user; %s", err.Error())
		return nil
	}

	return fmt.Errorf("Failed to create document on behalf of uphold user; status: %d", status)
}

// GetUser fetches the user for the given bearer token
func GetUser(token string) (*User, error) {
	var user *User
	var err error

	client, err := NewUpholdAPIClient(stringOrNil(token), stringOrNil("/v0/me"))
	if err != nil {
		return nil, err
	}

	status, err := client.Get("", nil, &user)
	if err != nil {
		log.Warningf("Failed to fetch uphold user on behalf of client id: %s; %s", upholdClientID, err.Error())
		return nil, err
	}

	if status == 200 {
		log.Debugf("Fetched uphold user %s on behalf of client id: %s; response: %s", *user.ID, upholdClientID)
		return user, nil
	}

	return nil, fmt.Errorf("Failed to fetch uphold user; status: %d", status)
}
