package indexer

import (
	"context"
	"net/http"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
)

const authHeader = "X-Indexer-API-Token"

type Client common.Client

// delete performs a DELETE request to the specific path against the server, assumes JSON response
func (c *Client) delete(ctx context.Context, response interface{}, path string, body interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Delete(ctx, response, path, body, headers)
}

// get performs a GET request to the specific path against the server, assumes JSON response
func (c *Client) get(ctx context.Context, response interface{}, path string, body interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Get(ctx, response, path, body, headers)
}

// getMsgpack performs a GET request to the specific path against the server, assumes msgpack response
func (c *Client) getMsgpack(ctx context.Context, response interface{}, path string, body interface{}, headers []*common.Header) error {
	return (*common.Client)(c).GetRawMsgpack(ctx, response, path, body, headers)
}

// getMsgpack performs a GET request to the specific path against the server, assumes msgpack response
func (c *Client) getRaw(ctx context.Context, path string, body interface{}, headers []*common.Header) ([]byte, error) {
	return (*common.Client)(c).GetRaw(ctx, path, body, headers)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (c *Client) post(ctx context.Context, response interface{}, path string, params interface{}, headers []*common.Header, body interface{}) error {
	return (*common.Client)(c).Post(ctx, response, path, params, headers, body)
}

// MakeClient is the factory for constructing a ClientV2 for a given endpoint.
func MakeClient(address string, apiToken string) (c *Client, err error) {
	commonClient, err := common.MakeClient(address, authHeader, apiToken)
	c = (*Client)(commonClient)
	return
}

// MakeClientWithHeaders is the factory for constructing a ClientV2 for a
// given endpoint with custom headers.
func MakeClientWithHeaders(address string, apiToken string, headers []*common.Header) (c *Client, err error) {
	commonClientWithHeaders, err := common.MakeClientWithHeaders(address, authHeader, apiToken, headers)
	c = (*Client)(commonClientWithHeaders)
	return
}

// MakeClientWithTransport is the factory for constructing a Client for a given endpoint with a
// custom HTTP Transport as well as optional additional user defined headers.
func MakeClientWithTransport(address string, apiToken string, headers []*common.Header, transport http.RoundTripper) (c *Client, err error) {
	commonClientWithTransport, err := common.MakeClientWithTransport(address, authHeader, apiToken, headers, transport)
	c = (*Client)(commonClientWithTransport)
	return
}

func (c *Client) HealthCheck() *HealthCheck {
	return &HealthCheck{c: c}
}

func (c *Client) SearchAccounts() *SearchAccounts {
	return &SearchAccounts{c: c}
}

func (c *Client) LookupAccountByID(accountId string) *LookupAccountByID {
	return &LookupAccountByID{c: c, accountId: accountId}
}

func (c *Client) LookupAccountAssets(accountId string) *LookupAccountAssets {
	return &LookupAccountAssets{c: c, accountId: accountId}
}

func (c *Client) LookupAccountCreatedAssets(accountId string) *LookupAccountCreatedAssets {
	return &LookupAccountCreatedAssets{c: c, accountId: accountId}
}

func (c *Client) LookupAccountAppLocalStates(accountId string) *LookupAccountAppLocalStates {
	return &LookupAccountAppLocalStates{c: c, accountId: accountId}
}

func (c *Client) LookupAccountCreatedApplications(accountId string) *LookupAccountCreatedApplications {
	return &LookupAccountCreatedApplications{c: c, accountId: accountId}
}

func (c *Client) LookupAccountTransactions(accountId string) *LookupAccountTransactions {
	return &LookupAccountTransactions{c: c, accountId: accountId}
}

func (c *Client) SearchForApplications() *SearchForApplications {
	return &SearchForApplications{c: c}
}

func (c *Client) LookupApplicationByID(applicationId uint64) *LookupApplicationByID {
	return &LookupApplicationByID{c: c, applicationId: applicationId}
}

func (c *Client) SearchForApplicationBoxes(applicationId uint64) *SearchForApplicationBoxes {
	return &SearchForApplicationBoxes{c: c, applicationId: applicationId}
}

func (c *Client) LookupApplicationBoxByIDAndName(applicationId uint64, name []byte) *LookupApplicationBoxByIDAndName {
	return (&LookupApplicationBoxByIDAndName{c: c, applicationId: applicationId}).name(name)
}

func (c *Client) LookupApplicationLogsByID(applicationId uint64) *LookupApplicationLogsByID {
	return &LookupApplicationLogsByID{c: c, applicationId: applicationId}
}

func (c *Client) SearchForAssets() *SearchForAssets {
	return &SearchForAssets{c: c}
}

func (c *Client) LookupAssetByID(assetId uint64) *LookupAssetByID {
	return &LookupAssetByID{c: c, assetId: assetId}
}

func (c *Client) LookupAssetBalances(assetId uint64) *LookupAssetBalances {
	return &LookupAssetBalances{c: c, assetId: assetId}
}

func (c *Client) LookupAssetTransactions(assetId uint64) *LookupAssetTransactions {
	return &LookupAssetTransactions{c: c, assetId: assetId}
}

func (c *Client) LookupBlock(roundNumber uint64) *LookupBlock {
	return &LookupBlock{c: c, roundNumber: roundNumber}
}

func (c *Client) LookupTransaction(txid string) *LookupTransaction {
	return &LookupTransaction{c: c, txid: txid}
}

func (c *Client) SearchForTransactions() *SearchForTransactions {
	return &SearchForTransactions{c: c}
}
