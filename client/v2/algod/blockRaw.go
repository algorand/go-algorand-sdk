package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// BlockRaw contains metadata required to execute a BlockRaw query.
type BlockRaw struct {
	c     *Client
	round uint64
	p     models.GetBlockParams
}

// Do executes the BlockRaw query and gets the results.
func (s *BlockRaw) Do(ctx context.Context, headers ...*common.Header) (result []byte, err error) {
	s.p.Format = "msgpack"
	return s.c.getRaw(ctx, fmt.Sprintf("/v2/blocks/%d", s.round), s.p, headers)
}
