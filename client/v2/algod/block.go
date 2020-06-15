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

type generatedBlockResponse struct {
	// Block header data.
	Block types.Block `json:"block"`

	// Optional certificate object. This is only included when the format is set to message pack.
	Cert *map[string]interface{} `json:"cert,omitempty"`
}

func (s *Block) Do(ctx context.Context, headers ...*common.Header) (result types.Block, err error) {
	s.p.Format = "msgpack"
	var response generatedBlockResponse
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/blocks/%d", s.round), s.p, headers)
	if err != nil {
		return
	}
	result = response.Block
	return
}
