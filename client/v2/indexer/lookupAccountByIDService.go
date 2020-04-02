package indexer

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type LookupAccountByIDService struct {
	c       *Client
	account string
	p       models.LookupAccountByIDParams
}

func (s *LookupAccountByIDService) Round(round uint64) *LookupAccountByIDService {
	s.p.Round = round
	return s
}

func (s *LookupAccountByIDService) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result models.Account, err error) {
	response := models.LookupAccountByIDResponse{}
	err = s.c.get(ctx, &response, fmt.Sprintf("/accounts/%s", s.account), s.p, headers)
	validRound = response.CurrentRound
	result = response.Account
	return
}
