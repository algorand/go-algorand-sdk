package common

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractError(t *testing.T) {
	testcases := []struct {
		name string
		code int
		err  error
	}{
		{name: "400", code: 400, err: BadRequest{HTTPError: HTTPError{StatusCode: 400, Message: ""}}},
		{name: "401", code: 401, err: InvalidToken{HTTPError: HTTPError{StatusCode: 401, Message: ""}}},
		{name: "404", code: 404, err: NotFound{HTTPError: HTTPError{StatusCode: 404, Message: ""}}},
		{name: "500", code: 500, err: InternalError{HTTPError: HTTPError{StatusCode: 500, Message: ""}}},
		{name: "418", code: 418, err: HTTPError{StatusCode: 418, Message: ""}},
		{name: "200", code: 200, err: nil},
		{name: "201", code: 201, err: nil},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := extractError(tc.code, []byte{})
			assert.Equal(t, tc.err, err)
			if err != nil {
				// Ensure concrete types are distinguishable at runtime.
				switch tc.code {
				case 400:
					_, ok := err.(BadRequest)
					assert.True(t, ok)
				case 401:
					_, ok := err.(InvalidToken)
					assert.True(t, ok)
				case 404:
					_, ok := err.(NotFound)
					assert.True(t, ok)
				case 500:
					_, ok := err.(InternalError)
					assert.True(t, ok)
				default:
					_, ok := err.(HTTPError)
					assert.True(t, ok)
				}
			}
		})
	}
}

func TestClient_Verbs(t *testing.T) {
	path := "/some/path"

	// Call each of the helper functions, they should make a request with the correct verb.
	testcases := []struct {
		expectedVerb string
		call         func(c *Client) error
	}{
		{
			expectedVerb: "GET",
			call: func(c *Client) error {
				return c.Get(context.Background(), nil, path, nil, nil)
			},
		},
		{
			expectedVerb: "POST",
			call: func(c *Client) error {
				return c.Post(context.Background(), nil, path, nil, nil, nil)
			},
		},
		{
			expectedVerb: "DELETE",
			call: func(c *Client) error {
				return c.Delete(context.Background(), nil, path, nil, nil)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.expectedVerb, func(t *testing.T) {
			var receivedMethod string
			var receivedPath string

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				receivedPath = r.URL.String()
			}))

			c, err := MakeClient(mockServer.URL, "API-Header", "ASDF")
			require.NoError(t, err)

			// Call the test function.
			err = tc.call(c)
			assert.Equal(t, tc.expectedVerb, receivedMethod)
			assert.Equal(t, path, receivedPath)
		})
	}
}

func TestClientWithTransport(t *testing.T) {
	var receivedMethod string
	var receivedPath string
	var receivedHeaderValue string
	path := "/some/path"

	const headerKey string = "hello"
	const headerValue string = "world"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.String()
		receivedHeaderValue = r.Header.Get(headerKey)
	}))

	var header []*Header = []*Header{{Key: headerKey, Value: headerValue}}
	var customTransport http.RoundTripper = &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c, err := MakeClientWithTransport(mockServer.URL, "API-Header", "ASDF", header, customTransport)
	require.NoError(t, err)

	// Call the test function.
	err = c.Get(context.Background(), nil, path, nil, nil)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, path, receivedPath)
	assert.Equal(t, headerValue, receivedHeaderValue)
	assert.Equal(t, c.transport, customTransport)
}
