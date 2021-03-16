package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type GetProof struct {
	c *Client

	round uint64

	txid string
}

func (s *GetProof) Do(ctx context.Context, headers ...*common.Header) (response models.ProofResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%v/transactions/%v/proof", s.round, s.txid), nil, headers)
	return
}
