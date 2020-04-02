package indexer

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type LookupBlockService struct {
	c     *Client
	round uint64
}

func (s *LookupBlockService) Do(ctx context.Context, headers ...*common.Header) (block models.Block, err error) {
	var response models.BlockResponse
	err = s.c.get(ctx, &response, fmt.Sprintf("/blocks/%d", s.round), nil, headers)
	block = models.Block(response)
	return
}
