package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
)

type TealDryrun struct {
	c *Client

	request models.DryrunRequest
}

func (s *TealDryrun) Do(ctx context.Context, headers ...*common.Header) (response models.DryrunResponse, err error) {
	err = s.c.post(ctx, &response, "/v2/teal/dryrun", msgpack.Encode(&s.request), headers)
	return
}
