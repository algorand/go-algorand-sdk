package common

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		tc := tc
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
