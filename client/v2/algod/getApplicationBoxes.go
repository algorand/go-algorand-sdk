package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// GetApplicationBoxesParams contains all of the query parameters for url serialization.
type GetApplicationBoxesParams struct {

	// Include include additional items in the response. Use `values` to include box
	// values. Multiple values can be comma-separated.
	Include []string `url:"include,omitempty,comma"`

	// Limit maximum number of boxes to return per page.
	Limit uint64 `url:"limit,omitempty"`

	// Max max number of box names to return. If max is not set, or max == 0, returns
	// all box-names.
	Max uint64 `url:"max,omitempty"`

	// Next a box name, in the goal app call arg form 'encoding:value', representing
	// the earliest box name to include in results. Use the next-token from a previous
	// response.
	Next string `url:"next,omitempty"`

	// Prefix a box name prefix, in the goal app call arg form 'encoding:value', to
	// filter results by. Only boxes whose names start with this prefix will be
	// returned.
	Prefix string `url:"prefix,omitempty"`

	// Round return box data from the given round. The round must be within the node's
	// available range.
	Round uint64 `url:"round,omitempty"`
}

// GetApplicationBoxes given an application ID, return all box names. No particular
// ordering is guaranteed. Request fails when client or server-side configured
// limits prevent returning all box names.
// Pagination mode is enabled when any of the following parameters are provided:
// limit, next, prefix, include, or round. In pagination mode box values can be
// requested and results are returned in sorted order.
// To paginate: use the next-token from a previous response as the next parameter
// in the following request. Pin the round parameter to the round value from the
// first page's response to ensure consistent results across pages. The server
// enforces a per-response byte limit, so fewer results than limit may be returned
// even when more exist; the presence of next-token is the only reliable signal
// that more data is available.
type GetApplicationBoxes struct {
	c *Client

	applicationId uint64

	p GetApplicationBoxesParams
}

// Include include additional items in the response. Use `values` to include box
// values. Multiple values can be comma-separated.
func (s *GetApplicationBoxes) Include(Include []string) *GetApplicationBoxes {
	s.p.Include = Include

	return s
}

// Limit maximum number of boxes to return per page.
func (s *GetApplicationBoxes) Limit(Limit uint64) *GetApplicationBoxes {
	s.p.Limit = Limit

	return s
}

// Max max number of box names to return. If max is not set, or max == 0, returns
// all box-names.
func (s *GetApplicationBoxes) Max(Max uint64) *GetApplicationBoxes {
	s.p.Max = Max

	return s
}

// Next a box name, in the goal app call arg form 'encoding:value', representing
// the earliest box name to include in results. Use the next-token from a previous
// response.
func (s *GetApplicationBoxes) Next(Next string) *GetApplicationBoxes {
	s.p.Next = Next

	return s
}

// Prefix a box name prefix, in the goal app call arg form 'encoding:value', to
// filter results by. Only boxes whose names start with this prefix will be
// returned.
func (s *GetApplicationBoxes) Prefix(Prefix string) *GetApplicationBoxes {
	s.p.Prefix = Prefix

	return s
}

// Round return box data from the given round. The round must be within the node's
// available range.
func (s *GetApplicationBoxes) Round(Round uint64) *GetApplicationBoxes {
	s.p.Round = Round

	return s
}

// Do performs the HTTP request
func (s *GetApplicationBoxes) Do(ctx context.Context, headers ...*common.Header) (response models.BoxesResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s/boxes", common.EscapeParams(s.applicationId)...), s.p, headers)
	return
}
