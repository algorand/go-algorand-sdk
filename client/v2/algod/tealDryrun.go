package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
)

// TealDryrun /v2/teal/dryrun
// Executes TEAL program(s) in context and returns debugging information about the
// execution.
type TealDryrun struct {
	c       *Client
	request models.DryrunRequest
}

// Do performs HTTP request
func (s *TealDryrun) Do(ctx context.Context,
	headers ...*common.Header) (response models.DryrunResponse, err error) {
	err = s.c.post(ctx, &response,
		"/v2/teal/dryrun", msgpack.Encode(&s.request), headers)
	return
}
