package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupApplicationLogsByIDParams contains all of the query parameters for url serialization.
type LookupApplicationLogsByIDParams struct {

	// Limit maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// MaxRound include results at or before the specified max-round.
	MaxRound uint64 `url:"max-round,omitempty"`

	// MinRound include results at or after the specified min-round.
	MinRound uint64 `url:"min-round,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`

	// SenderAddress only include transactions with this sender address.
	SenderAddress string `url:"sender-address,omitempty"`

	// Txid lookup the specific transaction by ID.
	Txid string `url:"txid,omitempty"`
}

// LookupApplicationLogsByID lookup application logs.
type LookupApplicationLogsByID struct {
	c *Client

	applicationId uint64

	p LookupApplicationLogsByIDParams
}

// Limit maximum number of results to return.
func (s *LookupApplicationLogsByID) Limit(Limit uint64) *LookupApplicationLogsByID {
	s.p.Limit = Limit
	return s
}

// MaxRound include results at or before the specified max-round.
func (s *LookupApplicationLogsByID) MaxRound(MaxRound uint64) *LookupApplicationLogsByID {
	s.p.MaxRound = MaxRound
	return s
}

// MinRound include results at or after the specified min-round.
func (s *LookupApplicationLogsByID) MinRound(MinRound uint64) *LookupApplicationLogsByID {
	s.p.MinRound = MinRound
	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *LookupApplicationLogsByID) Next(Next string) *LookupApplicationLogsByID {
	s.p.Next = Next
	return s
}

// SenderAddress only include transactions with this sender address.
func (s *LookupApplicationLogsByID) SenderAddress(SenderAddress string) *LookupApplicationLogsByID {
	s.p.SenderAddress = SenderAddress
	return s
}

// Txid lookup the specific transaction by ID.
func (s *LookupApplicationLogsByID) Txid(Txid string) *LookupApplicationLogsByID {
	s.p.Txid = Txid
	return s
}

// Do performs the HTTP request
func (s *LookupApplicationLogsByID) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationLogsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%v/logs", s.applicationId), s.p, headers)
	return
}
