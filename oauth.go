package uphold

import (
	"fmt"
)

// AuthorizeBearerToken synchronously authorizes a managed uphold API user using the environment-configured client id/secret and the given authorization code;
// note that it is the responsibility of the calling package to verify the provided state parameter, which should be a cryptographically secure random string
// used to protect against cross-site request forgery attacks. Packages which fail to verify the integrity of the state parameter provided alongside the code
// parameter passed into this function are vulnerable.
func AuthorizeBearerToken(code string) (*OAuthResponse, error) {
	var apiResponse *OAuthResponse
	var err error

	client, err := NewUpholdAPIClient(nil, nil)
	if err != nil {
		return nil, err
	}

	status, err := client.PostWWWFormURLEncoded("oauth2/token", map[string]interface{}{
		"code":       code,
		"grant_type": "authorization_code",
	}, &apiResponse)
	if err != nil {
		log.Warningf("Failed to authorize client credentials on behalf of client id: %s; %s", upholdClientID, err.Error())
		return nil, err
	}

	log.Debugf("Received %d status code in response to attempted client credentials authorization request on behalf of client id: %s; response: %s", status, upholdClientID, apiResponse)

	if status == 200 {
		log.Debugf("Resolved uphold %s access token: %s; refresh token: %s; scope: %s", apiResponse.TokenType, apiResponse.AccessToken, apiResponse.RefreshToken, apiResponse.Scope)
		// if response, responseOk := resp.(map[string]interface{}); responseOk {
		// 	apiResponse = &AccessTokenResponse{}
		// 	if accessToken, accessTokenOk := response["access_token"].(string); accessTokenOk {
		// 		apiResponse.AccessToken = accessToken
		// 	}
		// 	if tokenType, tokenTypeOk := response["token_type"].(string); tokenTypeOk {
		// 		apiResponse.TokenType = tokenType
		// 	}
		// 	if refreshToken, refreshTokenOk := response["refresh_token"].(string); refreshTokenOk {
		// 		apiResponse.RefreshToken = refreshToken
		// 	}
		// 	if scope, scopeOk := response["scope"].(string); scopeOk {
		// 		apiResponse.Scope = scope
		// 	}
		// 	log.Debugf("Resolved uphold %s access token: %s; refresh token: %s; scope: %s", apiResponse.TokenType, apiResponse.AccessToken, apiResponse.RefreshToken, apiResponse.Scope)
		// } else {
		// 	err = fmt.Errorf("Failed to parse client credentials API response on behalf of client id: %s; status code: %d", upholdClientID, status)
		// 	log.Warning(err.Error())
		// 	return nil, err
		// }
	} else {
		err = fmt.Errorf("Failed to authorize client credentials on behalf of client id: %s; status code: %d", upholdClientID, status)
		log.Warning(err.Error())
		return nil, err
	}

	return apiResponse, err
}

// AuthorizeClientCredentials synchronously authorizes an uphold API user using the environment-configured client id and secret
func AuthorizeClientCredentials(scope string) (*string, error) {
	var apiResponse *OAuthResponse
	var err error

	client, err := NewUpholdAPIClient(nil, nil)
	if err != nil {
		return nil, err
	}

	status, err := client.PostWWWFormURLEncoded("oauth2/token", map[string]interface{}{
		"grant_type": "client_credentials",
	}, &apiResponse)
	if err != nil {
		log.Warningf("Failed to authorize client credentials on behalf of client id: %s; %s", upholdClientID, err.Error())
		return nil, err
	}

	log.Debugf("Received %d status code in response to attempted client credentials authorization request on behalf of client id: %s; response: %s", status, upholdClientID, apiResponse)

	return apiResponse.AccessToken, err
}
