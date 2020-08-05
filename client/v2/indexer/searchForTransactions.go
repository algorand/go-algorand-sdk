package indexer

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

// SearchForTransactions /v2/transactions
// Search for transactions.
type SearchForTransactions struct {
	c *Client
	p models.SearchForTransactionsParams
}

// ApplicationId application ID
func (s *SearchForTransactions) ApplicationId(applicationId uint64) *SearchForTransactions {
	s.p.ApplicationId = applicationId
	return s
}

func (s *SearchForTransactions) NextToken(nextToken string) *SearchForTransactions {
	s.p.NextToken = nextToken
	return s
}

// NotePrefix specifies a prefix which must be contained in the note field.
func (s *SearchForTransactions) NotePrefix(prefix []byte) *SearchForTransactions {
	s.p.NotePrefix = base64.StdEncoding.EncodeToString(prefix)
	return s
}

func (s *SearchForTransactions) TxType(txtype string) *SearchForTransactions {
	s.p.TxType = txtype
	return s
}

// SigType filters just results using the specified type of signature:
// * sig - Standard
// * msig - MultiSig
// * lsig - LogicSig
func (s *SearchForTransactions) SigType(sigtype string) *SearchForTransactions {
	s.p.SigType = sigtype
	return s
}

// TXID lookup the specific transaction by ID.
func (s *SearchForTransactions) TXID(txid string) *SearchForTransactions {
	s.p.TxId = txid
	return s
}

// Round include results for the specified round.
func (s *SearchForTransactions) Round(rnd uint64) *SearchForTransactions {
	s.p.Round = rnd
	return s
}

// MinRound include results at or after the specified min-round.
func (s *SearchForTransactions) MinRound(min uint64) *SearchForTransactions {
	s.p.MinRound = min
	return s
}

// MaxRound include results at or before the specified max-round.
func (s *SearchForTransactions) MaxRound(max uint64) *SearchForTransactions {
	s.p.MaxRound = max
	return s
}

func (s *SearchForTransactions) AssetID(index uint64) *SearchForTransactions {
	s.p.AssetId = index
	return s
}

// Limit maximum number of results to return.
func (s *SearchForTransactions) Limit(limit uint64) *SearchForTransactions {
	s.p.Limit = limit
	return s
}

func (s *SearchForTransactions) BeforeTimeString(before string) *SearchForTransactions {
	s.p.BeforeTime = before
	return s
}

func (s *SearchForTransactions) AfterTimeString(after string) *SearchForTransactions {
	s.p.AfterTime = after
	return s
}

// BeforeTime include results before the given time. Must be an RFC 3339 formatted
// string.
func (s *SearchForTransactions) BeforeTime(before time.Time) *SearchForTransactions {
	beforeString := before.Format(time.RFC3339)
	return s.BeforeTimeString(beforeString)
}

// AfterTime include results after the given time. Must be an RFC 3339 formatted
// string.
func (s *SearchForTransactions) AfterTime(after time.Time) *SearchForTransactions {
	afterString := after.Format(time.RFC3339)
	return s.AfterTimeString(afterString)
}

// CurrencyGreaterThan results should have an amount greater than this value.
// MicroAlgos are the default currency unless an asset-id is provided, in which
// case the asset will be used.
func (s *SearchForTransactions) CurrencyGreaterThan(greaterThan uint64) *SearchForTransactions {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

// CurrencyLessThan results should have an amount less than this value. MicroAlgos
// are the default currency unless an asset-id is provided, in which case the asset
// will be used.
func (s *SearchForTransactions) CurrencyLessThan(lessThan uint64) *SearchForTransactions {
	s.p.CurrencyLessThan = lessThan
	return s
}

// AddressRole combine with the address parameter to define what type of address to
// search for.
func (s *SearchForTransactions) AddressRole(role string) *SearchForTransactions {
	s.p.AddressRole = role
	return s
}

func (s *SearchForTransactions) AddressString(address string) *SearchForTransactions {
	s.p.Address = address
	return s
}

// Address only include transactions with this address in one of the transaction
// fields.
func (s *SearchForTransactions) Address(address types.Address) *SearchForTransactions {
	return s.AddressString(address.String())
}

// ExcludeCloseTo combine with address and address-role parameters to define what
// type of address to search for. The close to fields are normally treated as a
// receiver, if you would like to exclude them set this parameter to true.
func (s *SearchForTransactions) ExcludeCloseTo(exclude bool) *SearchForTransactions {
	s.p.ExcludeCloseTo = exclude
	return s
}

// RekeyTo include results which include the rekey-to field.
func (s *SearchForTransactions) RekeyTo(rekeyTo bool) *SearchForTransactions {
	s.p.RekeyTo = rekeyTo
	return s
}

func (s *SearchForTransactions) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/transactions", s.p, headers)
	return
}
