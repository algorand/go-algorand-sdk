package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// AccountAssetsInformationParams contains all of the query parameters for url serialization.
type AccountAssetsInformationParams struct {

	// Limit maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`
}

// AccountAssetsInformation lookup an account's asset holdings.
type AccountAssetsInformation struct {
	c *Client

	address string

	p AccountAssetsInformationParams
}

// Limit maximum number of results to return.
func (s *AccountAssetsInformation) Limit(Limit uint64) *AccountAssetsInformation {
	s.p.Limit = Limit

	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *AccountAssetsInformation) Next(Next string) *AccountAssetsInformation {
	s.p.Next = Next

	return s
}

// Do performs the HTTP request
func (s *AccountAssetsInformation) Do(ctx context.Context, headers ...*common.Header) (response models.AccountAssetsInformationResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s/assets", common.EscapeParams(s.address)...), s.p, headers)
	return
}
