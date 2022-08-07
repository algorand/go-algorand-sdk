package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// TealDryRun /v2/teal/dryrun
type TealDryRun struct {
	c      *Client
	rawobj []byte
}

// Do performs HTTP request
func (s *TealDryRun) Do(
	ctx context.Context,
	headers ...*common.Header,
) (response models.DryrunResponse, err error) {
	err = s.c.post(ctx, &response, "/v2/teal/dryrun", nil, headers, s.rawobj)
	return
}
