package common

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractError(t *testing.T) {
	testcases := []struct {
		name    string
		code    int
		wantErr string
	}{
		{name: "400", code: 400, wantErr: "HTTP 400: "},
		{name: "401", code: 401, wantErr: "HTTP 401: "},
		{name: "404", code: 404, wantErr: "HTTP 404: "},
		{name: "500", code: 500, wantErr: "HTTP 500: "},
		{name: "200", code: 200},
		{name: "201", code: 201},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := extractError(tc.code, []byte{})
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			require.EqualError(t, err, tc.wantErr)

			switch tc.code {
			case 400:
				assert.IsType(t, BadRequest{}, err)
			case 401:
				assert.IsType(t, InvalidToken{}, err)
			case 404:
				assert.IsType(t, NotFound{}, err)
			case 500:
				assert.IsType(t, InternalError{}, err)
			}
		})
	}
}

func TestExtractErrorIncludesRESTData(t *testing.T) {
	err := extractError(400, []byte(`{"message":"txn dead","data":{"pool-error":"logic eval error","pc":7}}`))
	require.Error(t, err)
	require.EqualError(t, err, `HTTP 400: txn dead`)

	var httpErr *HTTPError
	require.True(t, errors.As(err, &httpErr))
	require.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
	assert.Equal(t, "txn dead", httpErr.Message)

	data := httpErr.Data
	require.NotNil(t, data)
	assert.Equal(t, "logic eval error", data["pool-error"])
	assert.Equal(t, float64(7), data["pc"])
}

func TestExtractErrorIgnoresInvalidJSON(t *testing.T) {
	err := extractError(400, []byte("not json"))
	require.Error(t, err)
	require.EqualError(t, err, "HTTP 400: not json")

	var httpErr *HTTPError
	require.True(t, errors.As(err, &httpErr))
	require.NotNil(t, httpErr)
	assert.Empty(t, httpErr.Message)
	assert.Nil(t, httpErr.Data)
	assert.Equal(t, []byte("not json"), httpErr.RawBody)
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
