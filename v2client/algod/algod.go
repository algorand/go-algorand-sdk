package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/v2client/common"
)

type ClientV2 common.Client

// get performs a GET request to the specific path against the server
func (client ClientV2) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return common.Client(client).Get(ctx, response, path, request, headers)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (client ClientV2) post(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return common.Client(client).Post(ctx, response, path, request, headers)
}

// MakeClient is the factory for constructing a ClientV2 for a given endpoint.
func MakeClient(address string, apiToken string) (c ClientV2, err error) {
	commonClient, err := common.MakeClient(address, apiToken)
	c = ClientV2(commonClient)
	return
}
