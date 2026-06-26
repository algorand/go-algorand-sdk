package algod

import (
	"context"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// SendRawTransactionParams contains all of the query parameters for url serialization.
type SendRawTransactionParams struct {

	// SkipPqAddressCheck skip post-quantum address checks, including the check that
	// rejects TEAL v13 or later LogicSig escrow addresses whose program hash is an
	// Edwards25519 curve point. This should only be used if you understand the risks
	// and know what you are doing.
	SkipPqAddressCheck bool `url:"skip-pq-address-check,omitempty"`
}

// SendRawTransaction broadcasts a raw transaction or transaction group to the
// network.
type SendRawTransaction struct {
	c *Client

	rawtxn []byte

	p SendRawTransactionParams
}

// SkipPqAddressCheck skip post-quantum address checks, including the check that
// rejects TEAL v13 or later LogicSig escrow addresses whose program hash is an
// Edwards25519 curve point. This should only be used if you understand the risks
// and know what you are doing.
func (s *SendRawTransaction) SkipPqAddressCheck(SkipPqAddressCheck bool) *SendRawTransaction {
	s.p.SkipPqAddressCheck = SkipPqAddressCheck

	return s
}

// Do performs the HTTP request
func (s *SendRawTransaction) Do(ctx context.Context, headers ...*common.Header) (txid string, err error) {
	var response models.PostTransactionsResponse
	// Set default Content-Type, if the user didn't specify it.
	addContentType := true
	for _, header := range headers {
		if strings.ToLower(header.Key) == "content-type" {
			addContentType = false
			break
		}
	}
	if addContentType {
		headers = append(headers, &common.Header{Key: "Content-Type", Value: "application/x-binary"})
	}
	err = s.c.post(ctx, &response, "/v2/transactions", nil, headers, s.rawtxn)
	txid = response.Txid
	return
}
