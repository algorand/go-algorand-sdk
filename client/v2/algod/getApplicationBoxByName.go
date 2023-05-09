package algod

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// GetApplicationBoxByNameParams contains all of the query parameters for url serialization.
type GetApplicationBoxByNameParams struct {

	// Name a box name, in the goal app call arg form 'encoding:value'. For ints, use
	// the form 'int:1234'. For raw bytes, use the form 'b64:A=='. For printable
	// strings, use the form 'str:hello'. For addresses, use the form 'addr:XYZ...'.
	Name string `url:"name,omitempty"`
}

// GetApplicationBoxByName given an application ID and box name, it returns the
// round, box name, and value (each base64 encoded). Box names must be in the goal
// app call arg encoding form 'encoding:value'. For ints, use the form 'int:1234'.
// For raw bytes, use the form 'b64:A=='. For printable strings, use the form
// 'str:hello'. For addresses, use the form 'addr:XYZ...'.
type GetApplicationBoxByName struct {
	c *Client

	applicationId uint64

	p GetApplicationBoxByNameParams
}

// name a box name, in the goal app call arg form 'encoding:value'. For ints, use
// the form 'int:1234'. For raw bytes, use the form 'b64:A=='. For printable
// strings, use the form 'str:hello'. For addresses, use the form 'addr:XYZ...'.
func (s *GetApplicationBoxByName) name(name []byte) *GetApplicationBoxByName {
	s.p.Name = "b64:" + base64.StdEncoding.EncodeToString(name)

	return s
}

// Do performs the HTTP request
func (s *GetApplicationBoxByName) Do(ctx context.Context, headers ...*common.Header) (response models.Box, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s/box", common.EscapeParams(s.applicationId)...), s.p, headers)
	return
}
