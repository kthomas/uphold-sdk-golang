package uphold

// AccessTokenResponse is the API response returned when an authorization code has been successfully upgraded to an access token
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
