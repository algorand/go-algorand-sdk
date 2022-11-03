package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// SearchForApplicationBoxesParams contains all of the query parameters for url serialization.
type SearchForApplicationBoxesParams struct {

	// Limit maximum number of results to return. There could be additional pages even
	// if the limit is not reached.
	Limit uint64 `url:"limit,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`
}

// SearchForApplicationBoxes given an application ID, returns the box names of that
// application sorted lexicographically.
type SearchForApplicationBoxes struct {
	c *Client

	applicationId uint64

	p SearchForApplicationBoxesParams
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *SearchForApplicationBoxes) Limit(Limit uint64) *SearchForApplicationBoxes {
	s.p.Limit = Limit

	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *SearchForApplicationBoxes) Next(Next string) *SearchForApplicationBoxes {
	s.p.Next = Next

	return s
}

// Do performs the HTTP request
func (s *SearchForApplicationBoxes) Do(ctx context.Context, headers ...*common.Header) (response models.BoxesResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s/boxes", common.EscapeParams(s.applicationId)...), s.p, headers)
	return
}
