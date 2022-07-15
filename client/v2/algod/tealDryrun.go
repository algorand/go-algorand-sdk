package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
)

// TealDryrun executes TEAL program(s) in context and returns debugging information
// about the execution. This endpoint is only enabled when a node's configuration
// file sets EnableDeveloperAPI to true.
type TealDryrun struct {
	c *Client

	request models.DryrunRequest
}

// Do performs the HTTP request
func (s *TealDryrun) Do(ctx context.Context, headers ...*common.Header) (response models.DryrunResponse, err error) {
	err = s.c.post(ctx, &response, "/v2/teal/dryrun", nil, headers, msgpack.Encode(&s.request))
	return
}
