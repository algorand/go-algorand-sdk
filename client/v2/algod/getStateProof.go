package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetStateProof get a state proof that covers a given round
type GetStateProof struct {
	c *Client

	round uint64
}

// Do performs the HTTP request
func (s *GetStateProof) Do(ctx context.Context, headers ...*common.Header) (response models.StateProof, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/stateproofs/%s", common.EscapeParams(s.round)...), nil, headers)
	return
}
