package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// rawRequestPaths is a set of paths where the body should not be urlencoded
var rawRequestPaths = map[string]bool{
	"/v2/transactions": true,
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
func extractError(resp *http.Response) error {
	if resp.StatusCode == 200 {
		return nil
	}

	errorBuf, _ := ioutil.ReadAll(resp.Body) // ignore returned error
	wrappedError := fmt.Errorf("HTTP %v: %s", resp.Status, errorBuf)
	switch code := resp.StatusCode; code {
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

// submitForm is a helper used for submitting (ex.) GETs and POSTs to the server
func (client *Client) submitFormRaw(ctx context.Context, path string, request interface{}, requestMethod string, encodeJSON bool, headers []*Header) (resp *http.Response, err error) {
	queryURL := client.serverURL
	queryURL.Path += path

	var req *http.Request
	var body io.Reader
	if request != nil {
		if requestMethod == "POST" && rawRequestPaths[path] {
			reqBytes, ok := request.([]byte)
			if !ok {
				return nil, fmt.Errorf("couldn't decode raw request as bytes")
			}
			body = bytes.NewBuffer(reqBytes)
		} else {
			v, err := query.Values(request)
			if err != nil {
				return nil, err
			}

			queryURL.RawQuery = mergeRawQueries(queryURL.RawQuery, v.Encode())
			if encodeJSON {
				jsonValue, _ := json.Marshal(request)
				body = bytes.NewBuffer(jsonValue)
			}
		}
	}

	req, err = http.NewRequest(requestMethod, queryURL.String(), body)
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
	err = extractError(resp)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	return resp, nil
}

func (client *Client) submitForm(ctx context.Context, response interface{}, path string, request interface{}, requestMethod string, encodeJSON bool, headers []*Header) error {
	resp, err := client.submitFormRaw(ctx, path, request, requestMethod, encodeJSON, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(&response)
}

// get performs a GET request to the specific path against the server
func (client *Client) Get(ctx context.Context, response interface{}, path string, request interface{}, headers []*Header) error {
	return client.submitForm(ctx, response, path, request, "GET", false /* encodeJSON */, headers)
}

func (client *Client) GetRawMsgpack(ctx context.Context, response interface{}, path string, request interface{}, headers []*Header) error {
	resp, err := client.submitFormRaw(ctx, path, request, "GET", false /* encodeJSON */, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dec := msgpack.NewDecoder(resp.Body)
	return dec.Decode(&response)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (client *Client) Post(ctx context.Context, response interface{}, path string, request interface{}, headers []*Header) error {
	return client.submitForm(ctx, response, path, request, "POST", true /* encodeJSON */, headers)
}
