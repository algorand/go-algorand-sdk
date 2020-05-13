package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type LookupAccountByID struct {
	c       *Client
	account string
	p       models.LookupAccountByIDParams
}

func (s *LookupAccountByID) Round(round uint64) *LookupAccountByID {
	s.p.Round = round
	return s
}

func (s *LookupAccountByID) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result models.Account, err error) {
	response := models.LookupAccountByIDResponse{}
	err = s.c.get(ctx, &response, fmt.Sprintf("%s/accounts/%s", versionPrefix, s.account), s.p, headers)
	validRound = response.CurrentRound
	result = response.Account
	return
}
