package indexer

import (
	"context"
	"time"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// SearchForBlockHeadersParams contains all of the query parameters for url serialization.
type SearchForBlockHeadersParams struct {

	// Absent accounts marked as absent in the block header's participation updates.
	// This parameter accepts a comma separated list of addresses.
	Absent []string `url:"absent,omitempty,comma"`

	// AfterTime include results after the given time. Must be an RFC 3339 formatted
	// string.
	AfterTime string `url:"after-time,omitempty"`

	// BeforeTime include results before the given time. Must be an RFC 3339 formatted
	// string.
	BeforeTime string `url:"before-time,omitempty"`

	// Expired accounts marked as expired in the block header's participation updates.
	// This parameter accepts a comma separated list of addresses.
	Expired []string `url:"expired,omitempty,comma"`

	// Limit maximum number of results to return. There could be additional pages even
	// if the limit is not reached.
	Limit uint64 `url:"limit,omitempty"`

	// MaxRound include results at or before the specified max-round.
	MaxRound uint64 `url:"max-round,omitempty"`

	// MinRound include results at or after the specified min-round.
	MinRound uint64 `url:"min-round,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`

	// Proposers accounts marked as proposer in the block header's participation
	// updates. This parameter accepts a comma separated list of addresses.
	Proposers []string `url:"proposers,omitempty,comma"`
}

// SearchForBlockHeaders search for block headers. Block headers are returned in
// ascending round order. Transactions are not included in the output.
type SearchForBlockHeaders struct {
	c *Client

	p SearchForBlockHeadersParams
}

// Absent accounts marked as absent in the block header's participation updates.
// This parameter accepts a comma separated list of addresses.
func (s *SearchForBlockHeaders) Absent(Absent []string) *SearchForBlockHeaders {
	s.p.Absent = Absent

	return s
}

// AfterTimeString include results after the given time. Must be an RFC 3339
// formatted string.
func (s *SearchForBlockHeaders) AfterTimeString(AfterTime string) *SearchForBlockHeaders {
	s.p.AfterTime = AfterTime

	return s
}

// AfterTime include results after the given time. Must be an RFC 3339 formatted
// string.
func (s *SearchForBlockHeaders) AfterTime(AfterTime time.Time) *SearchForBlockHeaders {
	AfterTimeStr := AfterTime.Format(time.RFC3339)

	return s.AfterTimeString(AfterTimeStr)
}

// BeforeTimeString include results before the given time. Must be an RFC 3339
// formatted string.
func (s *SearchForBlockHeaders) BeforeTimeString(BeforeTime string) *SearchForBlockHeaders {
	s.p.BeforeTime = BeforeTime

	return s
}

// BeforeTime include results before the given time. Must be an RFC 3339 formatted
// string.
func (s *SearchForBlockHeaders) BeforeTime(BeforeTime time.Time) *SearchForBlockHeaders {
	BeforeTimeStr := BeforeTime.Format(time.RFC3339)

	return s.BeforeTimeString(BeforeTimeStr)
}

// Expired accounts marked as expired in the block header's participation updates.
// This parameter accepts a comma separated list of addresses.
func (s *SearchForBlockHeaders) Expired(Expired []string) *SearchForBlockHeaders {
	s.p.Expired = Expired

	return s
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *SearchForBlockHeaders) Limit(Limit uint64) *SearchForBlockHeaders {
	s.p.Limit = Limit

	return s
}

// MaxRound include results at or before the specified max-round.
func (s *SearchForBlockHeaders) MaxRound(MaxRound uint64) *SearchForBlockHeaders {
	s.p.MaxRound = MaxRound

	return s
}

// MinRound include results at or after the specified min-round.
func (s *SearchForBlockHeaders) MinRound(MinRound uint64) *SearchForBlockHeaders {
	s.p.MinRound = MinRound

	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *SearchForBlockHeaders) Next(Next string) *SearchForBlockHeaders {
	s.p.Next = Next

	return s
}

// Proposers accounts marked as proposer in the block header's participation
// updates. This parameter accepts a comma separated list of addresses.
func (s *SearchForBlockHeaders) Proposers(Proposers []string) *SearchForBlockHeaders {
	s.p.Proposers = Proposers

	return s
}

// Do performs the HTTP request
func (s *SearchForBlockHeaders) Do(ctx context.Context, headers ...*common.Header) (response models.BlockHeadersResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/block-headers", s.p, headers)
	return
}
