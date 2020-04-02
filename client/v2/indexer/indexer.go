package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

type Client common.Client

// get performs a GET request to the specific path against the server
func (c *Client) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Get(ctx, response, path, request, headers)
}

// MakeClient is the factory for constructing an IndexerClient for a given endpoint.
func MakeClient(address string, apiToken string) (c *Client, err error) {
	commonClient, err := common.MakeClient(address, apiToken)
	c = (*Client)(commonClient)
	return
}

func (c *Client) NewLookupAssetBalancesService(index uint64) *LookupAssetBalancesService {
	return &LookupAssetBalancesService{c: c, index: index}
}

func (c *Client) NewLookupAssetTransactionsService(index uint64) *LookupAssetTransactionsService {
	return &LookupAssetTransactionsService{c: c, index: index}
}

func (c *Client) NewLookupAccountTransactionsService(account string) *LookupAccountTransactionsService {
	return &LookupAccountTransactionsService{c: c, account: account}
}

func (c *Client) NewLookupBlockService(round uint64) *LookupBlockService {
	return &LookupBlockService{c: c, round: round}
}

func (c *Client) NewLookupAccountByIDService(account string) *LookupAccountByIDService {
	return &LookupAccountByIDService{c: c, account: account}
}

func (c *Client) NewLookupAssetByIDService(index uint64) *LookupAssetByIDService {
	return &LookupAssetByIDService{c: c, index: index}
}

func (c *Client) NewSearchAccountsService() *SearchAccountsService {
	return &SearchAccountsService{c: c}
}

func (c *Client) NewSearchForTransactionsService() *SearchForTransactionsService {
	return &SearchForTransactionsService{c: c}
}

func (c *Client) NewSearchForAssetsService() *SearchForAssetsService {
	return &SearchForAssetsService{c: c}
}
