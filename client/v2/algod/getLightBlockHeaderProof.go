package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetLightBlockHeaderProof gets a proof for a given light block header inside a
// state proof commitment
type GetLightBlockHeaderProof struct {
	c *Client

	round uint64
}

// Do performs the HTTP request
func (s *GetLightBlockHeaderProof) Do(ctx context.Context, headers ...*common.Header) (response models.LightBlockHeaderProof, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%s/lightheader/proof", common.EscapeParams(s.round)...), nil, headers)
	return
}
