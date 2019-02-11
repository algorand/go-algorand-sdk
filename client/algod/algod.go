package algod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	authHeader           = "X-Algo-API-Token"
	healthCheckEndpoint  = "/health"
	apiVersionPathPrefix = "/v1"
)

// unversionedPaths ais a set of paths that should not be prefixed by the API version
var unversionedPaths = map[string]bool{
	"/versions": true,
	"/health":   true,
}

// rawRequestPaths is a set of paths where the body should not be urlencoded
var rawRequestPaths = map[string]bool{
	"/transactions": true,
}

// Client manages the REST interface for a calling user.
type Client struct {
	serverURL url.URL
	apiToken  string
}

// MakeClient is the factory for constructing a Client for a given endpoint
func MakeClient(address string, apiToken string) (c Client, err error) {
	url, err := url.Parse(address)
	if err != nil {
		return
	}
	c = Client{
		serverURL: *url,
		apiToken:  apiToken,
	}
	return
}

// extractError checks if the response signifies an error (for now, StatusCode != 200).
// If so, it returns the error.
// Otherwise, it returns nil.
func extractError(resp *http.Response) error {
	if resp.StatusCode == 200 {
		return nil
	}

	errorBuf, _ := ioutil.ReadAll(resp.Body) // ignore returned error
	return fmt.Errorf("HTTP %v: %s", resp.Status, errorBuf)
}

// stripTransaction gets a transaction of the form "tx-XXXXXXXX" and truncates the "tx-" part, if it starts with "tx-"
func stripTransaction(tx string) string {
	if strings.HasPrefix(tx, "tx-") {
		return strings.SplitAfter(tx, "-")[1]
	}
	return tx
}

// submitForm is a helper used for submitting (ex.) GETs and POSTs to the server
func (client Client) submitForm(response interface{}, path string, request interface{}, requestMethod string, encodeJSON bool) error {
	var err error
	queryURL := client.serverURL
	queryURL.Path = path

	// Handle version prefix
	if !unversionedPaths[path] {
		queryURL.Path = strings.Join([]string{apiVersionPathPrefix, path}, "")
	}

	var req *http.Request
	var body io.Reader

	if request != nil {
		if rawRequestPaths[path] {
			reqBytes, ok := request.([]byte)
			if !ok {
				return fmt.Errorf("couldn't decode raw request as bytes")
			}
			body = bytes.NewBuffer(reqBytes)
		} else {
			v, err := query.Values(request)
			if err != nil {
				return err
			}

			queryURL.RawQuery = v.Encode()
			if encodeJSON {
				jsonValue, _ := json.Marshal(request)
				body = bytes.NewBuffer(jsonValue)
			}
		}
	}

	req, err = http.NewRequest(requestMethod, queryURL.String(), body)
	if err != nil {
		return err
	}

	// If we add another endpoint that does not require auth, we should add a
	// requiresAuth argument to submitForm rather than checking here
	if path != healthCheckEndpoint {
		req.Header.Set(authHeader, client.apiToken)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = extractError(resp)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(&response)
}

// get performs a GET request to the specific path against the server
func (client Client) get(response interface{}, path string, request interface{}) error {
	return client.submitForm(response, path, request, "GET", false /* encodeJSON */)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (client Client) post(response interface{}, path string, request interface{}) error {
	return client.submitForm(response, path, request, "POST", true /* encodeJSON */)
}

// as post, but with MethodPut
func (client Client) put(response interface{}, path string, request interface{}) error {
	return client.submitForm(response, path, request, "PUT", true /* encodeJSON */)
}

// as post, but with MethodPatch
func (client Client) patch(response interface{}, path string, request interface{}) error {
	return client.submitForm(response, path, request, "PATCH", true /* encodeJSON */)
}
