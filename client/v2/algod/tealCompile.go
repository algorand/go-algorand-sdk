package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type TealCompile struct {
	c *Client

	source []byte
}

func (s *TealCompile) Do(ctx context.Context, headers ...*common.Header) (response models.CompileResponse, err error) {
	err = s.c.post(ctx, &response, "/v2/teal/compile", s.source, headers)
	return
}
