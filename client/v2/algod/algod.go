package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

type Client common.Client

// get performs a GET request to the specific path against the server
func (client Client) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return common.Client(client).Get(ctx, response, path, request, headers)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (client Client) post(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return common.Client(client).Post(ctx, response, path, request, headers)
}

// MakeClient is the factory for constructing a ClientV2 for a given endpoint.
func MakeClient(address string, apiToken string) (c Client, err error) {
	commonClient, err := common.MakeClient(address, apiToken)
	c = Client(commonClient)
	return
}
