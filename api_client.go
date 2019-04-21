package uphold

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultContentType = "application/json"

// APIClient is a generic base class for calling the uphold API
type APIClient struct {
	Host     string
	Path     string
	Scheme   string
	Token    *string
	Username *string
	Password *string
}

// NewUpholdAPIClient initializes an APIClient using the environment-configured client id and secret
// to construct an HTTP basic authorization header, unless a non-nil bearer access token is provided.
func NewUpholdAPIClient(token, baseURI *string) (*APIClient, error) {
	apiURL, err := url.Parse(upholdAPIBaseURL)
	if err != nil {
		log.Warningf("Failed to parse uphold API base url; %s", err.Error())
		return nil, err
	}

	path := ""
	if baseURI != nil {
		path = *baseURI
	}

	var client *APIClient

	if token != nil {
		client = &APIClient{
			Host:   apiURL.Host,
			Scheme: apiURL.Scheme,
			Path:   path,
			Token:  token,
		}
	} else {
		client = &APIClient{
			Host:     apiURL.Host,
			Scheme:   apiURL.Scheme,
			Path:     path,
			Username: stringOrNil(upholdClientID),
			Password: stringOrNil(upholdClientSecret),
		}
	}

	return client, nil
}

// NewUnauthorizedAPIClient initializes an APIClient without API credentials
func NewUnauthorizedAPIClient(baseURI *string) (*APIClient, error) {
	apiURL, err := url.Parse(upholdAPIBaseURL)
	if err != nil {
		log.Warningf("Failed to parse uphold API base url; %s", err.Error())
		return nil, err
	}

	path := ""
	if baseURI != nil {
		path = *baseURI
	}

	return &APIClient{
		Host:   apiURL.Host,
		Scheme: apiURL.Scheme,
		Path:   path,
		Token:  nil,
	}, nil
}

func (c *APIClient) sendRequest(method, urlString, contentType string, params map[string]interface{}, response interface{}) (status int, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		Timeout: time.Second * 30,
	}

	mthd := strings.ToUpper(method)
	reqURL, err := url.Parse(urlString)
	if err != nil {
		log.Warningf("Failed to parse URL for uphold API (%s %s) invocation; %s", method, urlString, err.Error())
		return -1, err
	}

	if mthd == "GET" && params != nil {
		q := reqURL.Query()
		for name := range params {
			if val, valOk := params[name].(string); valOk {
				q.Set(name, val)
			}
		}
		reqURL.RawQuery = q.Encode()
	}

	headers := map[string][]string{
		"Accept-Encoding": {"gzip, deflate"},
		"Accept-Language": {"en-us"},
		"Accept":          {"application/json"},
	}
	if c.Username != nil && c.Password != nil {
		headers["Authorization"] = []string{buildBasicAuthorizationHeader(*c.Username, *c.Password)}
	} else if c.Token != nil {
		headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", *c.Token)}
	}

	var req *http.Request

	if mthd == "POST" || mthd == "PUT" {
		var payload []byte
		if contentType == "application/json" {
			payload, err = json.Marshal(params)
			if err != nil {
				log.Warningf("Failed to marshal JSON payload for uphold API (%s %s) invocation; %s", method, urlString, err.Error())
				return -1, err
			}
		} else if contentType == "application/x-www-form-urlencoded" {
			urlEncodedForm := url.Values{}
			for key, val := range params {
				if valStr, valOk := val.(string); valOk {
					urlEncodedForm.Add(key, valStr)
				} else {
					log.Warningf("Failed to marshal application/x-www-form-urlencoded parameter: %s; value was non-string", key)
				}
			}
			payload = []byte(urlEncodedForm.Encode())
		}

		req, _ = http.NewRequest(method, urlString, bytes.NewReader(payload))
		headers["Content-Type"] = []string{contentType}
	} else {
		req = &http.Request{
			URL:    reqURL,
			Method: mthd,
		}
	}

	req.Header = headers

	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Warningf("Failed to invoke uphold API (%s %s) method: %s; %s", method, urlString, err.Error())
		return 0, err
	}

	log.Debugf("Received %v response for uphold API (%s %s) invocation", resp.StatusCode, method, urlString)

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("Failed to unmarshal uphold API (%s %s) response: %s; %s", method, urlString, buf.Bytes(), err.Error())
	}

	log.Debugf("Invocation of uphold API (%s %s) succeeded (%v-byte response)", method, urlString, buf.Len())
	return resp.StatusCode, nil
}

// Get constructs and synchronously sends an API GET request
func (c *APIClient) Get(uri string, params map[string]interface{}, response interface{}) (status int, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("GET", url, defaultContentType, params, response)
}

// Post constructs and synchronously sends an API POST request
func (c *APIClient) Post(uri string, params map[string]interface{}, response interface{}) (status int, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("POST", url, defaultContentType, params, response)
}

// PostWWWFormURLEncoded constructs and synchronously sends an API POST request using
func (c *APIClient) PostWWWFormURLEncoded(uri string, params map[string]interface{}, response interface{}) (status int, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("POST", url, "application/x-www-form-urlencoded", params, response)
}

// Put constructs and synchronously sends an API PUT request
func (c *APIClient) Put(uri string, params map[string]interface{}, response interface{}) (status int, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("PUT", url, defaultContentType, params, response)
}

// Delete constructs and synchronously sends an API DELETE request
func (c *APIClient) Delete(uri string) (status int, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("DELETE", url, defaultContentType, nil, nil)
}

func (c *APIClient) buildURL(uri string) string {
	path := c.Path
	if len(path) == 1 && path == "/" {
		path = ""
	} else if len(path) > 1 && strings.Index(path, "/") != 0 {
		path = fmt.Sprintf("/%s", path)
	}
	return fmt.Sprintf("%s://%s%s/%s", c.Scheme, c.Host, path, uri)
}

func buildBasicAuthorizationHeader(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
}
