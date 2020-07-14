package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// TealDryrun /v2/teal/dryrun
// Executes TEAL program(s) in context and returns debugging information about the
// execution.
type TealDryrun struct {
	c       *Client
	request DryrunRequest
}

// Do performs HTTP request
func (s *TealDryrun) Do(ctx context.Context,
	headers ...*common.Header) (response models.DryrunResponse, err error) {
	err = s.c.post(ctx, &response,
		"/v2/teal/dryrun", s.request, headers)
	return
}
