package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// AccountApplicationsInformationParams contains all of the query parameters for url serialization.
type AccountApplicationsInformationParams struct {

	// Include additional items in the response.
	Include []string `url:"include,omitempty,comma"`

	// Limit maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`
}

// AccountApplicationsInformation lookup an account's application holdings (local
// state and params if the account is the creator).
type AccountApplicationsInformation struct {
	c *Client

	address string

	p AccountApplicationsInformationParams
}

// Include additional items in the response.
func (s *AccountApplicationsInformation) Include(Include []string) *AccountApplicationsInformation {
	s.p.Include = Include

	return s
}

// Limit maximum number of results to return.
func (s *AccountApplicationsInformation) Limit(Limit uint64) *AccountApplicationsInformation {
	s.p.Limit = Limit

	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *AccountApplicationsInformation) Next(Next string) *AccountApplicationsInformation {
	s.p.Next = Next

	return s
}

// Do performs the HTTP request
func (s *AccountApplicationsInformation) Do(ctx context.Context, headers ...*common.Header) (response models.AccountApplicationsInformationResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s/applications", common.EscapeParams(s.address)...), s.p, headers)
	return
}
