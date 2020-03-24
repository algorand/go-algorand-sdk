package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/v2client/common"
	"net/url"
)

type IndexerClient common.Client

// get performs a GET request to the specific path against the server
func (client IndexerClient) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*Header) error {
	return common.Client(client).Get(ctx, response, path, request, headers)
}

// MakeClient is the factory for constructing an IndexerClient for a given endpoint.
func MakeClient(address string, apiToken string) (c IndexerClient, err error) {
	commonClient, err := common.MakeClient(address, apiToken)
	c = IndexerClient(commonClient)
	return
}
