package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupApplicationBoxByIDandNameParams contains all of the query parameters for url serialization.
type LookupApplicationBoxByIDandNameParams struct {

	// Name a box name in goal-arg form 'encoding:value'. For ints, use the form
	// 'int:1234'. For raw bytes, use the form 'b64:A=='. For printable strings, use
	// the form 'str:hello'. For addresses, use the form 'addr:XYZ...'.
	Name string `url:"name,omitempty"`
}

// LookupApplicationBoxByIDandName given an application ID and box name, returns
// base64 encoded box name and value. Box names must be in the goal-arg form
// 'encoding:value'. For ints, use the form 'int:1234'. For raw bytes, encode base
// 64 and use 'b64' prefix as in 'b64:A=='. For printable strings, use the form
// 'str:hello'. For addresses, use the form 'addr:XYZ...'.
type LookupApplicationBoxByIDandName struct {
	c *Client

	applicationId uint64

	p LookupApplicationBoxByIDandNameParams
}

// Name a box name in goal-arg form 'encoding:value'. For ints, use the form
// 'int:1234'. For raw bytes, use the form 'b64:A=='. For printable strings, use
// the form 'str:hello'. For addresses, use the form 'addr:XYZ...'.
func (s *LookupApplicationBoxByIDandName) Name(Name string) *LookupApplicationBoxByIDandName {
	s.p.Name = Name
	return s
}

// Do performs the HTTP request
func (s *LookupApplicationBoxByIDandName) Do(ctx context.Context, headers ...*common.Header) (response models.Box, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s/box", common.EscapeParams(s.applicationId)...), s.p, headers)
	return
}
