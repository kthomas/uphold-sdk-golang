package uphold

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// APIClient is a generic base class for calling the uphold API
type APIClient struct {
	Host     string
	Path     string
	Scheme   string
	Token    *string
	Username *string
	Password *string
}

func (c *APIClient) sendRequest(method, urlString string, params map[string]interface{}) (status int, response interface{}, err error) {
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
		return -1, nil, err
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
		payload, err := json.Marshal(params)
		if err != nil {
			log.Warningf("Failed to marshal JSON payload for uphold API (%s %s) invocation; %s", method, urlString, err.Error())
			return -1, nil, err
		}
		req, _ = http.NewRequest(method, urlString, bytes.NewReader(payload))
		headers["Content-Type"] = []string{"application/json"}
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
		return 0, nil, err
	}

	log.Debugf("Received %v response for uphold API (%s %s) invocation", resp.StatusCode, method, urlString)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("Failed to unmarshal uphold API (%s %s) response: %s; %s", method, urlString, buf.Bytes(), err.Error())
	}

	log.Debugf("Invocation of uphold API (%s %s) succeeded (%v-byte response)", method, urlString, buf.Len())
	return resp.StatusCode, response, nil
}

// Get constructs and synchronously sends an API GET request
func (c *APIClient) Get(uri string, params map[string]interface{}) (status int, response interface{}, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("GET", url, params)
}

// Post constructs and synchronously sends an API POST request
func (c *APIClient) Post(uri string, params map[string]interface{}) (status int, response interface{}, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("POST", url, params)
}

// Put constructs and synchronously sends an API PUT request
func (c *APIClient) Put(uri string, params map[string]interface{}) (status int, response interface{}, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("PUT", url, params)
}

// Delete constructs and synchronously sends an API DELETE request
func (c *APIClient) Delete(uri string) (status int, response interface{}, err error) {
	url := c.buildURL(uri)
	return c.sendRequest("DELETE", url, nil)
}

func buildBasicAuthorizationHeader(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *APIClient) buildURL(uri string) string {
	return fmt.Sprintf("%s://%s/%s/%s", c.Scheme, c.Host, c.Path, uri)
}
