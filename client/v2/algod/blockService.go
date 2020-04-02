package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type BlockService struct {
	c     *Client
	round uint64
}

func (s *BlockService) Do(ctx context.Context, headers ...*common.Header) (result types.Block, err error) {
	response := models.GetBlockResponse{}
	err = s.c.get(ctx, &response, fmt.Sprintf("/blocks/%d", s.round), models.NewBlockParams(), headers)
	if err != nil {
		return
	}
	err = result.FromBase64String(response.Blockb64)
	return
}
