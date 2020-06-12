package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

/**
 * /v2/accounts/{account-id}
 * Lookup account information.
 */
type LookupAccountByID struct {
	c         *Client
	p         models.LookupAccountByIDParams
	accountId string
}

/**
 * Include results for the specified round.
 */
func (s *LookupAccountByID) Round(round uint64) *LookupAccountByID {
	s.p.Round = round
	return s
}

func (s *LookupAccountByID) Do(ctx context.Context,
	headers ...*common.Header) (validRound uint64, result models.Account, err error) {
	response := models.LookupAccountByIDResponse{}
	err = s.c.get(ctx, &response,
		fmt.Sprintf("/v2/accounts/%s", s.accountId), s.p, headers)
	validRound = response.CurrentRound
	result = response.Account
	return
}
