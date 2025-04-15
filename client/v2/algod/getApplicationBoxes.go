package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// GetApplicationBoxesParams contains all of the query parameters for url serialization.
type GetApplicationBoxesParams struct {

	// Max maximum number of boxes to return. Server may impose a lower limit.
	Max uint64 `url:"max,omitempty"`

	// Next a box name, in the goal app call arg form 'encoding:value'. When provided,
	// the returned boxes begin (lexographically) with the supplied name. Callers may
	// implement pagination by reinvoking the endpoint with the token from a previous
	// call's next-token.
	Next string `url:"next,omitempty"`

	// Prefix a box name prefix, in the goal app call arg form 'encoding:value'. For
	// ints, use the form 'int:1234'. For raw bytes, use the form 'b64:A=='. For
	// printable strings, use the form 'str:hello'. For addresses, use the form
	// 'addr:XYZ...'.
	Prefix string `url:"prefix,omitempty"`

	// Values if true, box values will be returned.
	Values bool `url:"values,omitempty"`
}

// GetApplicationBoxes given an application ID, return boxes in lexographical order
// by name. If the results must be truncated, a next-token is supplied to continue
// the request.
type GetApplicationBoxes struct {
	c *Client

	applicationId uint64

	p GetApplicationBoxesParams
}

// Max maximum number of boxes to return. Server may impose a lower limit.
func (s *GetApplicationBoxes) Max(Max uint64) *GetApplicationBoxes {
	s.p.Max = Max

	return s
}

// Next a box name, in the goal app call arg form 'encoding:value'. When provided,
// the returned boxes begin (lexographically) with the supplied name. Callers may
// implement pagination by reinvoking the endpoint with the token from a previous
// call's next-token.
func (s *GetApplicationBoxes) Next(Next string) *GetApplicationBoxes {
	s.p.Next = Next

	return s
}

// Prefix a box name prefix, in the goal app call arg form 'encoding:value'. For
// ints, use the form 'int:1234'. For raw bytes, use the form 'b64:A=='. For
// printable strings, use the form 'str:hello'. For addresses, use the form
// 'addr:XYZ...'.
func (s *GetApplicationBoxes) Prefix(Prefix string) *GetApplicationBoxes {
	s.p.Prefix = Prefix

	return s
}

// Values if true, box values will be returned.
func (s *GetApplicationBoxes) Values(Values bool) *GetApplicationBoxes {
	s.p.Values = Values

	return s
}

// Do performs the HTTP request
func (s *GetApplicationBoxes) Do(ctx context.Context, headers ...*common.Header) (response models.BoxesResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s/boxes", common.EscapeParams(s.applicationId)...), s.p, headers)
	return
}
