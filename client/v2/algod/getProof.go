package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetProofParams contains all of the query parameters for url serialization.
type GetProofParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`

	// Hashtype the type of hash function used to create the proof, must be one of:
	// * sha512_256
	// * sha256
	Hashtype string `url:"hashtype,omitempty"`
}

// GetProof get a Merkle proof for a transaction in a block.
type GetProof struct {
	c *Client

	round uint64
	txid  string

	p GetProofParams
}

// Hashtype the type of hash function used to create the proof, must be one of:
// * sha512_256
// * sha256
func (s *GetProof) Hashtype(Hashtype string) *GetProof {
	s.p.Hashtype = Hashtype
	return s
}

// Do performs the HTTP request
func (s *GetProof) Do(ctx context.Context, headers ...*common.Header) (response models.ProofResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%v/transactions/%v/proof", s.round, s.txid), s.p, headers)
	return
}
