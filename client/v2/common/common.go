package common

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/algorand/go-algorand-sdk/v2/encoding/json"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/google/go-querystring/query"
)

// rawRequestPaths is a set of paths where the body should not be urlencoded
var rawRequestPaths = map[string]bool{
	"/v2/transactions":     true,
	"/v2/teal/compile":     true,
	"/v2/teal/disassemble": true,
	"/v2/teal/dryrun":      true,
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
	transport http.RoundTripper
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

// MakeClientWithTransport is the factory for constructing a Client for a given endpoint with a
// custom HTTP Transport as well as optional additional user defined headers.
func MakeClientWithTransport(address string, apiHeader, apiToken string, headers []*Header, transport http.RoundTripper) (c *Client, err error) {
	c, err = MakeClientWithHeaders(address, apiHeader, apiToken, headers)
	if err != nil {
		return
	}

	c.transport = transport

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
	if code >= 200 && code < 300 {
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
func (client *Client) submitFormRaw(ctx context.Context, path string, params interface{}, requestMethod string, encodeJSON bool, headers []*Header, body interface{}) (resp *http.Response, err error) {
	queryURL := client.serverURL
	queryURL.Path += path

	var (
		req        *http.Request
		bodyReader io.Reader
		v          url.Values
	)

	if params != nil {
		v, err = query.Values(params)
		if err != nil {
			return nil, err
		}
	}

	if requestMethod == "POST" && rawRequestPaths[path] {
		reqBytes, ok := body.([]byte)
		if !ok {
			return nil, fmt.Errorf("couldn't decode raw body as bytes")
		}
		bodyReader = bytes.NewBuffer(reqBytes)
	} else if encodeJSON {
		jsonValue := json.Encode(params)
		bodyReader = bytes.NewBuffer(jsonValue)
	}

	queryURL.RawQuery = mergeRawQueries(queryURL.RawQuery, v.Encode())

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

	httpClient := &http.Client{Transport: client.transport}
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

func (client *Client) submitForm(ctx context.Context, response interface{}, path string, params interface{}, requestMethod string, encodeJSON bool, headers []*Header, body interface{}) error {
	resp, err := client.submitFormRaw(ctx, path, params, requestMethod, encodeJSON, headers, body)
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
		return responseErr
	}

	// Attempt to unmarshal a response regardless of whether or not there was an error.
	err = json.LenientDecode(bodyBytes, response)
	if responseErr != nil {
		// Even if there was an unmarshalling error, return the HTTP error first if there was one.
		return responseErr
	}
	return err
}

// Delete performs a DELETE request to the specific path against the server
func (client *Client) Delete(ctx context.Context, response interface{}, path string, params interface{}, headers []*Header) error {
	return client.submitForm(ctx, response, path, params, "DELETE", false /* encodeJSON */, headers, nil)
}

// Get performs a GET request to the specific path against the server
func (client *Client) Get(ctx context.Context, response interface{}, path string, params interface{}, headers []*Header) error {
	return client.submitForm(ctx, response, path, params, "GET", false /* encodeJSON */, headers, nil)
}

// GetRaw performs a GET request to the specific path against the server and returns the raw body bytes.
func (client *Client) GetRaw(ctx context.Context, path string, params interface{}, headers []*Header) (response []byte, err error) {
	var resp *http.Response
	resp, err = client.submitFormRaw(ctx, path, params, "GET", false /* encodeJSON */, headers, nil)
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
func (client *Client) GetRawMsgpack(ctx context.Context, response interface{}, path string, params interface{}, headers []*Header) error {
	resp, err := client.submitFormRaw(ctx, path, params, "GET", false /* encodeJSON */, headers, nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var bodyBytes []byte
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %+v", err)
		}

		return extractError(resp.StatusCode, bodyBytes)
	}

	dec := msgpack.NewLenientDecoder(resp.Body)
	return dec.Decode(&response)
}

// Post sends a POST request to the given path with the given body object.
// No query parameters will be sent if body is nil.
// response must be a pointer to an object as post writes the response there.
func (client *Client) Post(ctx context.Context, response interface{}, path string, params interface{}, headers []*Header, body interface{}) error {
	return client.submitForm(ctx, response, path, params, "POST", true /* encodeJSON */, headers, body)
}

// Helper function for correctly formatting and escaping URL path parameters.
// Used in the generated API client code.
func EscapeParams(params ...interface{}) []interface{} {
	paramsStr := make([]interface{}, len(params))
	for i, param := range params {
		switch v := param.(type) {
		case string:
			paramsStr[i] = url.PathEscape(v)
		default:
			paramsStr[i] = fmt.Sprintf("%v", v)
		}
	}

	return paramsStr
}
