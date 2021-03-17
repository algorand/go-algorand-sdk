package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type GetProofParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`
}

type GetProof struct {
	c *Client

	round uint64

	txid string

	p GetProofParams
}

func (s *GetProof) Do(ctx context.Context, headers ...*common.Header) (response models.ProofResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%v/transactions/%v/proof", s.round, s.txid), s.p, headers)
	return
}
