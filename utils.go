package uphold

import (
	"fmt"
	"net/url"
)

// WebAuthorizationURL returns the webapp authorization URL for the given scope
func WebAuthorizationURL(scope string) string {
	return fmt.Sprintf("%s/authorize/%s?scope=%s", upholdBaseURL, upholdClientID, url.QueryEscape(scope))
}

// WebAuthorizationAllScopesURL returns the webapp authorization URL requesting all supported scopes
func WebAuthorizationAllScopesURL() string {
	return WebAuthorizationURL(upholdSupportedScopes)
}

func stringOrNil(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
