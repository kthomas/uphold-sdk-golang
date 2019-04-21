package uphold

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

	log.Debugf("Received %d status code when attempting to craete uphold user: %s; response: %s", status, upholdClientID, user)

	return user, err
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
		log.Warningf("Failed to fetch user on behalf of client id: %s; %s", upholdClientID, err.Error())
		return nil, err
	}

	log.Debugf("Received %d status code when attempting to fetch user on behalf of client id: %s; response: %s", status, upholdClientID, user)

	return user, err
}
