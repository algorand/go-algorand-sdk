package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
)

// SetBlockTimeStampOffset sets the timestamp offset (seconds) for blocks in dev
// mode. Providing an offset of 0 will unset this value and try to use the real
// clock for the timestamp.
type SetBlockTimeStampOffset struct {
	c *Client

	offset uint64
}

// Do performs the HTTP request
func (s *SetBlockTimeStampOffset) Do(ctx context.Context, headers ...*common.Header) (response string, err error) {
	err = s.c.post(ctx, &response, fmt.Sprintf("/v2/devmode/blocks/offset/%s", common.EscapeParams(s.offset)...), nil, headers, nil)
	return
}
