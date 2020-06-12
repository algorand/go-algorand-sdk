package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

/**
 * /v2/blocks/{round-number}
 * Lookup block.
 */
type LookupBlock struct {
	c           *Client
	roundNumber uint64
}

func (s *LookupBlock) Do(ctx context.Context,
	headers ...*common.Header) (response models.Block, err error) {
	err = s.c.get(ctx, &response,
		fmt.Sprintf("/v2/blocks/%d", s.roundNumber), nil, headers)
	return
}
