package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// LookupBlockParams contains all of the query parameters for url serialization.
type LookupBlockParams struct {

	// HeaderOnly header only flag. When this is set to true, returned block does not
	// contain the transactions
	HeaderOnly bool `url:"header-only,omitempty"`
}

// LookupBlock lookup block.
type LookupBlock struct {
	c *Client

	roundNumber uint64

	p LookupBlockParams
}

// HeaderOnly header only flag. When this is set to true, returned block does not
// contain the transactions
func (s *LookupBlock) HeaderOnly(HeaderOnly bool) *LookupBlock {
	s.p.HeaderOnly = HeaderOnly

	return s
}

// Do performs the HTTP request
func (s *LookupBlock) Do(ctx context.Context, headers ...*common.Header) (response models.Block, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%s", common.EscapeParams(s.roundNumber)...), s.p, headers)
	return
}
