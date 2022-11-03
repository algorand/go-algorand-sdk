package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetTransactionProofParams contains all of the query parameters for url serialization.
type GetTransactionProofParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`

	// Hashtype the type of hash function used to create the proof, must be one of:
	// * sha512_256
	// * sha256
	Hashtype string `url:"hashtype,omitempty"`
}

// GetTransactionProof get a proof for a transaction in a block.
type GetTransactionProof struct {
	c *Client

	round uint64
	txid  string

	p GetTransactionProofParams
}

// Hashtype the type of hash function used to create the proof, must be one of:
// * sha512_256
// * sha256
func (s *GetTransactionProof) Hashtype(Hashtype string) *GetTransactionProof {
	s.p.Hashtype = Hashtype

	return s
}

// Do performs the HTTP request
func (s *GetTransactionProof) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionProofResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%s/transactions/%s/proof", common.EscapeParams(s.round, s.txid)...), s.p, headers)
	return
}
