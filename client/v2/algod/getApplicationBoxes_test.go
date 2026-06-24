package algod

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetApplicationBoxes_QuerySerialization(t *testing.T) {
	testcases := []struct {
		name     string
		setup    func(*GetApplicationBoxes)
		expected map[string]string
	}{
		{
			name: "default query params",
			setup: func(g *GetApplicationBoxes) {
				// No setup, use defaults
			},
			expected: map[string]string{},
		},
		{
			name: "max only",
			setup: func(g *GetApplicationBoxes) {
				g.Max(100)
			},
			expected: map[string]string{
				"max": "100",
			},
		},
		{
			name: "pagination with limit",
			setup: func(g *GetApplicationBoxes) {
				g.Limit(20)
			},
			expected: map[string]string{
				"limit": "20",
			},
		},
		{
			name: "pagination with next token",
			setup: func(g *GetApplicationBoxes) {
				g.Next("b64:AAECAw==")
			},
			expected: map[string]string{
				"next": "b64:AAECAw==",
			},
		},
		{
			name: "prefix filter",
			setup: func(g *GetApplicationBoxes) {
				g.Prefix("b64:AAECAw==")
			},
			expected: map[string]string{
				"prefix": "b64:AAECAw==",
			},
		},
		{
			name: "include values",
			setup: func(g *GetApplicationBoxes) {
				g.Include([]string{"values"})
			},
			expected: map[string]string{
				"include": "values",
			},
		},
		{
			name: "round filter",
			setup: func(g *GetApplicationBoxes) {
				g.Round(98765)
			},
			expected: map[string]string{
				"round": "98765",
			},
		},
		{
			name: "full pagination query",
			setup: func(g *GetApplicationBoxes) {
				g.Next("b64:AAECAw==").
					Prefix("b64:AAECAw==").
					Include([]string{"values"}).
					Limit(20).
					Round(98765)
			},
			expected: map[string]string{
				"next":    "b64:AAECAw==",
				"prefix":  "b64:AAECAw==",
				"include": "values",
				"limit":   "20",
				"round":   "98765",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var receivedQuery map[string][]string

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedQuery = r.URL.Query()
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(models.BoxesResponse{})
			}))
			defer mockServer.Close()

			algodClient, err := MakeClient(mockServer.URL, "test-token")
			require.NoError(t, err)

			g := algodClient.GetApplicationBoxes(1234)

			tc.setup(g)

			// Execute the request
			_, _ = g.Do(context.Background())

			// Verify query params
			for key, expectedValue := range tc.expected {
				values, ok := receivedQuery[key]
				require.True(t, ok, "expected query param %s to be present", key)
				assert.Equal(t, expectedValue, values[0], "query param %s mismatch", key)
			}
		})
	}
}

func TestGetApplicationBoxes_ResponseDecoding(t *testing.T) {
	testcases := []struct {
		name         string
		responseBody string
		expected     models.BoxesResponse
	}{
		{
			name:         "empty response",
			responseBody: `{"boxes":[],"application-id":1234}`,
			expected: models.BoxesResponse{
				ApplicationId: 1234,
				Boxes:         []models.BoxDescriptor{},
			},
		},
		{
			name:         "boxes with names only",
			responseBody: `{"boxes":[{"name":"AAECAw=="}],"application-id":1234}`,
			expected: models.BoxesResponse{
				ApplicationId: 1234,
				Boxes: []models.BoxDescriptor{
					{Name: []byte{0, 1, 2, 3}},
				},
			},
		},
		{
			name:         "boxes with names and values",
			responseBody: `{"boxes":[{"name":"AAECAw==","value":"BAUG"}],"application-id":1234}`,
			expected: models.BoxesResponse{
				ApplicationId: 1234,
				Boxes: []models.BoxDescriptor{
					{Name: []byte{0, 1, 2, 3}, Value: []byte{4, 5, 6}},
				},
			},
		},
		{
			name:         "response with next-token",
			responseBody: `{"boxes":[],"application-id":1234,"next-token":"b64:Bwg="}`,
			expected: models.BoxesResponse{
				ApplicationId: 1234,
				Boxes:         []models.BoxDescriptor{},
				NextToken:     "b64:Bwg=",
			},
		},
		{
			name:         "response with round",
			responseBody: `{"boxes":[],"application-id":1234,"round":98765}`,
			expected: models.BoxesResponse{
				ApplicationId: 1234,
				Boxes:         []models.BoxDescriptor{},
				Round:         98765,
			},
		},
		{
			name:         "full response",
			responseBody: `{"boxes":[{"name":"AAECAw==","value":"BAUG"}],"application-id":1234,"next-token":"b64:Bwg=","round":11}`,
			expected: models.BoxesResponse{
				ApplicationId: 1234,
				Boxes: []models.BoxDescriptor{
					{Name: []byte{0, 1, 2, 3}, Value: []byte{4, 5, 6}},
				},
				NextToken: "b64:Bwg=",
				Round:     11,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tc.responseBody))
			}))
			defer mockServer.Close()

			algodClient, err := MakeClient(mockServer.URL, "test-token")
			require.NoError(t, err)

			g := algodClient.GetApplicationBoxes(1234)

			resp, err := g.Do(context.Background())
			require.NoError(t, err)

			assert.Equal(t, tc.expected.ApplicationId, resp.ApplicationId)
			assert.Equal(t, tc.expected.NextToken, resp.NextToken)
			assert.Equal(t, tc.expected.Round, resp.Round)
			assert.Equal(t, len(tc.expected.Boxes), len(resp.Boxes))

			for i, expectedBox := range tc.expected.Boxes {
				assert.Equal(t, expectedBox.Name, resp.Boxes[i].Name)
				assert.Equal(t, expectedBox.Value, resp.Boxes[i].Value)
			}
		})
	}
}
