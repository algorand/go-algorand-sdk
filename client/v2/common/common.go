package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/google/go-querystring/query"
)

// rawRequestPaths is a set of paths where the body should not be urlencoded
var rawRequestPaths = map[string]bool{
	"/v2/transactions": true,
	"/v2/teal/compile": true,
	"/v2/teal/dryrun":  true,
}

// Header is a struct for custom headers.
type Header struct {
	Key   string
	Value string
}

// Client manages the REST interface for a calling user.
type Client struct {
	serverURL url.URL
	apiHeader string
	apiToken  string
	headers   []*Header
}

// MakeClient is the factory for constructing a Client for a given endpoint.
func MakeClient(address string, apiHeader, apiToken string) (c *Client, err error) {
	url, err := url.Parse(address)
	if err != nil {
		return
	}

	c = &Client{
		serverURL: *url,
		apiHeader: apiHeader,
		apiToken:  apiToken,
	}
	return
}

// MakeClientWithHeaders is the factory for constructing a Client for a given endpoint with additional user defined headers.
func MakeClientWithHeaders(address string, apiHeader, apiToken string, headers []*Header) (c *Client, err error) {
	c, err = MakeClient(address, apiHeader, apiToken)
	if err != nil {
		return
	}

	c.headers = append(c.headers, headers...)

	return
}

type BadRequest error
type InvalidToken error
type NotFound error
type InternalError error

// extractError checks if the response signifies an error.
// If so, it returns the error.
// Otherwise, it returns nil.
func extractError(code int, errorBuf []byte) error {
	if code == 200 {
		return nil
	}

	wrappedError := fmt.Errorf("HTTP %v: %s", code, errorBuf)
	switch code {
	case 400:
		return BadRequest(wrappedError)
	case 401:
		return InvalidToken(wrappedError)
	case 404:
		return NotFound(wrappedError)
	case 500:
		return InternalError(wrappedError)
	default:
		return wrappedError
	}
}

// mergeRawQueries merges two raw queries, appending an "&" if both are non-empty
func mergeRawQueries(q1, q2 string) string {
	if q1 == "" {
		return q2
	}
	if q2 == "" {
		return q1
	}
	return q1 + "&" + q2
}

// submitFormRaw is a helper used for submitting (ex.) GETs and POSTs to the server
func (client *Client) submitFormRaw(ctx context.Context, path string, body interface{}, requestMethod string, encodeJSON bool, headers []*Header) (resp *http.Response, err error) {
	queryURL := client.serverURL
	queryURL.Path += path

	var req *http.Request
	var bodyReader io.Reader
	if body != nil {
		if requestMethod == "POST" && rawRequestPaths[path] {
			reqBytes, ok := body.([]byte)
			if !ok {
				return nil, fmt.Errorf("couldn't decode raw body as bytes")
			}
			bodyReader = bytes.NewBuffer(reqBytes)
		} else {
			v, err := query.Values(body)
			if err != nil {
				return nil, err
			}

			queryURL.RawQuery = mergeRawQueries(queryURL.RawQuery, v.Encode())
			if encodeJSON {
				jsonValue, _ := json.Marshal(body)
				bodyReader = bytes.NewBuffer(jsonValue)
			}
		}
	}

	req, err = http.NewRequest(requestMethod, queryURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	// Supply the client token.
	req.Header.Set(client.apiHeader, client.apiToken)
	// Add the client headers.
	for _, header := range client.headers {
		req.Header.Add(header.Key, header.Value)
	}
	// Add the request headers.
	for _, header := range headers {
		req.Header.Add(header.Key, header.Value)
	}

	httpClient := &http.Client{}
	req = req.WithContext(ctx)
	resp, err = httpClient.Do(req)

	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	return resp, nil
}

func (client *Client) submitForm(ctx context.Context, response interface{}, path string, body interface{}, requestMethod string, encodeJSON bool, headers []*Header) error {
	resp, err := client.submitFormRaw(ctx, path, body, requestMethod, encodeJSON, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	responseErr := extractError(resp.StatusCode, bodyBytes)

	// The caller wants a string
	if strResponse, ok := response.(*string); ok {
		*strResponse = string(bodyBytes)
		return err
	}

	// Attempt to unmarshal a response regardless of whether or not there was an error.
	err = json.Unmarshal(bodyBytes, response)
	if responseErr != nil {
		// Even if there was an unmarshalling error, return the HTTP error first if there was one.
		return responseErr
	}
	return err
}

// Get performs a GET request to the specific path against the server
func (client *Client) Get(ctx context.Context, response interface{}, path string, body interface{}, headers []*Header) error {
	return client.submitForm(ctx, response, path, body, "GET", false /* encodeJSON */, headers)
}

// GetRaw performs a GET request to the specific path against the server and returns the raw body bytes.
func (client *Client) GetRaw(ctx context.Context, path string, body interface{}, headers []*Header) (response []byte, err error) {
	var resp *http.Response
	resp, err = client.submitFormRaw(ctx, path, body, "GET", false /* encodeJSON */, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, extractError(resp.StatusCode, bodyBytes)
}

// GetRawMsgpack performs a GET request to the specific path against the server and returns the decoded messagepack response.
func (client *Client) GetRawMsgpack(ctx context.Context, response interface{}, path string, body interface{}, headers []*Header) error {
	resp, err := client.submitFormRaw(ctx, path, body, "GET", false /* encodeJSON */, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	dec := msgpack.NewDecoder(resp.Body)
	return dec.Decode(&response)
}

// Post sends a POST request to the given path with the given body object.
// No query parameters will be sent if body is nil.
// response must be a pointer to an object as post writes the response there.
func (client *Client) Post(ctx context.Context, response interface{}, path string, body interface{}, headers []*Header) error {
	return client.submitForm(ctx, response, path, body, "POST", true /* encodeJSON */, headers)
}
