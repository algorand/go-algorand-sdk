package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// BlockParams contains all of the query parameters for url serialization.
type BlockParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded. If
	// not provided, defaults to JSON.
	Format string `url:"format,omitempty"`

	// HeaderOnly if true, only the block header (exclusive of payset or certificate)
	// may be included in response.
	HeaderOnly bool `url:"header-only,omitempty"`
}

// Block get the block for the given round.
type Block struct {
	c *Client

	round uint64

	p BlockParams
}

// HeaderOnly if true, only the block header (exclusive of payset or certificate)
// may be included in response.
func (s *Block) HeaderOnly(HeaderOnly bool) *Block {
	s.p.HeaderOnly = HeaderOnly

	return s
}

// Do performs the HTTP request
func (s *Block) Do(ctx context.Context, headers ...*common.Header) (result types.Block, err error) {
	var response models.BlockResponse

	s.p.Format = "msgpack"
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/blocks/%d", s.round), s.p, headers)
	if err != nil {
		return
	}

	result = response.Block
	return
}
