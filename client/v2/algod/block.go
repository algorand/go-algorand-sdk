package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type Block struct {
	c     *Client
	round uint64
	p     models.GetBlockParams
}

type blockResponse struct {
	block types.Block `codec:"block"`
}

func (s *Block) Do(ctx context.Context, headers ...*common.Header) (result types.Block, err error) {
	s.p.Format = "msgpack"
	var response blockResponse
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/blocks/%d", s.round), s.p, headers)
	if err != nil {
		return
	}
	result = response.block
	return
}
