package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

const indexerAuthHeader = "X-Indexer-API-Token"

type Client common.Client

// get performs a GET request to the specific path against the server
func (c *Client) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Get(ctx, response, path, request, headers)
}

// MakeClient is the factory for constructing an IndexerClient for a given endpoint.
func MakeClient(address string, apiToken string) (c *Client, err error) {
	commonClient, err := common.MakeClient(address, indexerAuthHeader, apiToken)
	c = (*Client)(commonClient)
	return
}

func (c *Client) HealthCheck() *HealthCheck {
	return &HealthCheck{c: c}
}

func (c *Client) LookupAssetBalances(index uint64) *LookupAssetBalances {
	return &LookupAssetBalances{c: c, index: index}
}

func (c *Client) LookupAssetTransactions(index uint64) *LookupAssetTransactions {
	return &LookupAssetTransactions{c: c, index: index}
}

func (c *Client) LookupAccountTransactions(account string) *LookupAccountTransactions {
	return &LookupAccountTransactions{c: c, account: account}
}

func (c *Client) LookupBlock(round uint64) *LookupBlock {
	return &LookupBlock{c: c, round: round}
}

func (c *Client) LookupAccountByID(account string) *LookupAccountByID {
	return &LookupAccountByID{c: c, account: account}
}

func (c *Client) LookupAssetByID(index uint64) *LookupAssetByID {
	return &LookupAssetByID{c: c, index: index}
}

func (c *Client) SearchAccounts() *SearchAccounts {
	return &SearchAccounts{c: c}
}

func (c *Client) SearchForTransactions() *SearchForTransactions {
	return &SearchForTransactions{c: c}
}

func (c *Client) SearchForAssets() *SearchForAssets {
	return &SearchForAssets{c: c}
}
