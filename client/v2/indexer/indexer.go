package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

type Client common.Client

// get performs a GET request to the specific path against the server
func (client Client) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return common.Client(client).Get(ctx, response, path, request, headers)
}

// MakeClient is the factory for constructing an IndexerClient for a given endpoint.
func MakeClient(address string, apiToken string) (c Client, err error) {
	commonClient, err := common.MakeClient(address, apiToken)
	c = Client(commonClient)
	return
}
