package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type AccountInformationService struct {
	c       *Client
	account string
}

func (s *AccountInformationService) AccountInformation(ctx context.Context, headers ...*common.Header) (result models.Account, err error) {
	err = s.c.get(ctx, &result, fmt.Sprintf("/accounts/%s", s.account), nil, headers)
	return
}
