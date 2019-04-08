package uphold

// GetUser fetches the user for the given bearer token
func GetUser(token string) (*User, error) {
	var user *User
	var err error

	client, err := NewUpholdAPIClient(stringOrNil(token), stringOrNil("/v0/me"))
	if err != nil {
		return nil, err
	}

	status, err := client.Get("/", nil, &user)
	if err != nil {
		log.Warningf("Failed to fetch user on behalf of client id: %s; %s", upholdClientID, err.Error())
		return nil, err
	}

	log.Debugf("Received %d status code when attempting to fetch user on behalf of client id: %s; response: %s", status, upholdClientID, user)

	return user, err
}
