package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupApplicationBoxByIDAndNameParams contains all of the query parameters for url serialization.
type LookupApplicationBoxByIDAndNameParams struct {

	// Name a box name in goal-arg form 'encoding:value'. For ints, use the form
	// 'int:1234'. For raw bytes, use the form 'b64:A=='. For printable strings, use
	// the form 'str:hello'. For addresses, use the form 'addr:XYZ...'.
	Name string `url:"name,omitempty"`
}

// LookupApplicationBoxByIDAndName given an application ID and box name, returns
// base64 encoded box name and value. Box names must be in the goal app call arg
// form 'encoding:value'. For ints, use the form 'int:1234'. For raw bytes, encode
// base 64 and use 'b64' prefix as in 'b64:A=='. For printable strings, use the
// form 'str:hello'. For addresses, use the form 'addr:XYZ...'.
type LookupApplicationBoxByIDAndName struct {
	c *Client

	applicationId uint64

	p LookupApplicationBoxByIDAndNameParams
}

// Name a box name in goal-arg form 'encoding:value'. For ints, use the form
// 'int:1234'. For raw bytes, use the form 'b64:A=='. For printable strings, use
// the form 'str:hello'. For addresses, use the form 'addr:XYZ...'.
func (s *LookupApplicationBoxByIDAndName) Name(Name string) *LookupApplicationBoxByIDAndName {
	s.p.Name = Name
	return s
}

// Do performs the HTTP request
func (s *LookupApplicationBoxByIDAndName) Do(ctx context.Context, headers ...*common.Header) (response models.Box, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s/box", common.EscapeParams(s.applicationId)...), s.p, headers)
	return
}
