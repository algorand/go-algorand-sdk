package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Header is a struct for custom headers.
// TODO (ejr) is the usage of custom headers necessary for this?
type Header struct {
	Key   string
	Value string
}

//TODO: should algod.go be made into a common used by indexer? or is it custom enough that it needs to be written separately?

const (
	authHeader           = "X-Indexer-API-Token" // TODO EJR support new token structure
	healthCheckEndpoint  = "/health"             // TODO EJR indexer has different auth scheme
	apiVersionPathPrefix = "/v2"                 // TODO EJR support v1 and v2; indexer does not have version prefix I don't think
)

// unversionedPaths ais a set of paths that should not be prefixed by the API version
var unversionedPaths = map[string]bool{
	//none right now
}

// rawRequestPaths is a set of paths where the body should not be urlencoded
// TODO (ejr) these need to be updated
var rawRequestPaths = map[string]bool{
	"/transactions": true,
}

// Client manages the REST interface for a calling user.
type Client struct {
	serverURL url.URL
	apiToken  string
	headers   []*Header
}

// MakeClient is the factory for constructing a Client for a given endpoint.
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

// MakeClientWithHeaders is the factory for constructing a Client for a given endpoint with additional user defined headers.
func MakeClientWithHeaders(address string, apiToken string, headers []*Header) (c Client, err error) {
	c, err = MakeClient(address, apiToken)
	if err != nil {
		return
	}

	c.headers = append(c.headers, headers...)

	return
}

// extractError checks if the response signifies an error (for now, StatusCode != 200).
// If so, it returns the error.
// Otherwise, it returns nil.
// TODO EJR new errors can be handled here likely
func extractError(resp *http.Response) error {
	if resp.StatusCode == 200 {
		return nil
	}

	errorBuf, _ := ioutil.ReadAll(resp.Body) // ignore returned error
	return fmt.Errorf("HTTP %v: %s", resp.Status, errorBuf)
}

// mergeRawQueries merges two raw queries, appending an "&" if both are non-empty
func mergeRawQueries(q1, q2 string) string {
	if q1 == "" {
		return q2
	} else if q2 == "" {
		return q1
	} else {
		return q1 + "&" + q2
	}
}

// submitForm is a helper used for submitting (ex.) GETs and POSTs to the server
func (client Client) submitFormRaw(ctx context.Context, path string, request interface{}, requestMethod string, encodeJSON bool, headers []*Header) (resp *http.Response, err error) {
	queryURL := client.serverURL

	// Handle version prefix
	if !unversionedPaths[path] {
		queryURL.Path += strings.Join([]string{apiVersionPathPrefix, path}, "")
	} else {
		queryURL.Path += path
	}

	var req *http.Request
	var body io.Reader

	if request != nil {
		if rawRequestPaths[path] {
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

	// If we add another endpoint that does not require auth, we should add a
	// requiresAuth argument to submitForm rather than checking here
	// TODO requiresAuth will need to be used for algodclient v2 shutdown and keyreg
	if path != healthCheckEndpoint {
		req.Header.Set(authHeader, client.apiToken)
	}
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

func (client Client) submitForm(ctx context.Context, response interface{}, path string, request interface{}, requestMethod string, encodeJSON bool, headers []*Header) error {
	resp, err := client.submitFormRaw(ctx, path, request, requestMethod, encodeJSON, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(&response)
}

// get performs a GET request to the specific path against the server
func (client Client) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*Header) error {
	return client.submitForm(ctx, response, path, request, "GET", false /* encodeJSON */, headers)
}
