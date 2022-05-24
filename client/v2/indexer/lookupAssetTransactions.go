package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAssetTransactionsParams contains all of the query parameters for url serialization.
type LookupAssetTransactionsParams struct {

	// AddressString only include transactions with this address in one of the
	// transaction fields.
	AddressString string `url:"address,omitempty"`

	// AddressRole combine with the address parameter to define what type of address to
	// search for.
	AddressRole string `url:"address-role,omitempty"`

	// AfterTime include results after the given time. Must be an RFC 3339 formatted
	// string.
	AfterTime string `url:"after-time,omitempty"`

	// BeforeTime include results before the given time. Must be an RFC 3339 formatted
	// string.
	BeforeTime string `url:"before-time,omitempty"`

	// CurrencyGreaterThan results should have an amount greater than this value.
	// MicroAlgos are the default currency unless an asset-id is provided, in which
	// case the asset will be used.
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	// CurrencyLessThan results should have an amount less than this value. MicroAlgos
	// are the default currency unless an asset-id is provided, in which case the asset
	// will be used.
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	// ExcludeCloseTo combine with address and address-role parameters to define what
	// type of address to search for. The close to fields are normally treated as a
	// receiver, if you would like to exclude them set this parameter to true.
	ExcludeCloseTo bool `url:"exclude-close-to,omitempty"`

	// Limit maximum number of results to return. There could be additional pages even
	// if the limit is not reached.
	Limit uint64 `url:"limit,omitempty"`

	// MaxRound include results at or before the specified max-round.
	MaxRound uint64 `url:"max-round,omitempty"`

	// MinRound include results at or after the specified min-round.
	MinRound uint64 `url:"min-round,omitempty"`

	// NextToken the next page of results. Use the next token provided by the previous
	// results.
	NextToken string `url:"next,omitempty"`

	// NotePrefix specifies a prefix which must be contained in the note field.
	NotePrefix string `url:"note-prefix,omitempty"`

	// RekeyTo include results which include the rekey-to field.
	RekeyTo bool `url:"rekey-to,omitempty"`

	// Round include results for the specified round.
	Round uint64 `url:"round,omitempty"`

	// SigType sigType filters just results using the specified type of signature:
	// * sig - Standard
	// * msig - MultiSig
	// * lsig - LogicSig
	SigType string `url:"sig-type,omitempty"`

	// TxType
	TxType string `url:"tx-type,omitempty"`

	// TXID lookup the specific transaction by ID.
	TXID string `url:"txid,omitempty"`
}

// LookupAssetTransactions lookup transactions for an asset. Transactions are
// returned oldest to newest.
type LookupAssetTransactions struct {
	c *Client

	assetId uint64

	p LookupAssetTransactionsParams
}

// AddressString only include transactions with this address in one of the
// transaction fields.
func (s *LookupAssetTransactions) AddressString(AddressString string) *LookupAssetTransactions {
	s.p.AddressString = AddressString
	return s
}

// AddressRole combine with the address parameter to define what type of address to
// search for.
func (s *LookupAssetTransactions) AddressRole(AddressRole string) *LookupAssetTransactions {
	s.p.AddressRole = AddressRole
	return s
}

// AfterTimeString include results after the given time. Must be an RFC 3339
// formatted string.
func (s *LookupAssetTransactions) AfterTimeString(AfterTime string) *LookupAssetTransactions {
	s.p.AfterTime = AfterTime
	return s
}

// AfterTime include results after the given time. Must be an RFC 3339 formatted
// string.
func (s *LookupAssetTransactions) AfterTime(AfterTime time.Time) *LookupAssetTransactions {
	AfterTimeStr := AfterTime.Format(time.RFC3339)

	return s.AfterTimeString(AfterTimeStr)
}

// BeforeTimeString include results before the given time. Must be an RFC 3339
// formatted string.
func (s *LookupAssetTransactions) BeforeTimeString(BeforeTime string) *LookupAssetTransactions {
	s.p.BeforeTime = BeforeTime
	return s
}

// BeforeTime include results before the given time. Must be an RFC 3339 formatted
// string.
func (s *LookupAssetTransactions) BeforeTime(BeforeTime time.Time) *LookupAssetTransactions {
	BeforeTimeStr := BeforeTime.Format(time.RFC3339)

	return s.BeforeTimeString(BeforeTimeStr)
}

// CurrencyGreaterThan results should have an amount greater than this value.
// MicroAlgos are the default currency unless an asset-id is provided, in which
// case the asset will be used.
func (s *LookupAssetTransactions) CurrencyGreaterThan(CurrencyGreaterThan uint64) *LookupAssetTransactions {
	s.p.CurrencyGreaterThan = CurrencyGreaterThan
	return s
}

// CurrencyLessThan results should have an amount less than this value. MicroAlgos
// are the default currency unless an asset-id is provided, in which case the asset
// will be used.
func (s *LookupAssetTransactions) CurrencyLessThan(CurrencyLessThan uint64) *LookupAssetTransactions {
	s.p.CurrencyLessThan = CurrencyLessThan
	return s
}

// ExcludeCloseTo combine with address and address-role parameters to define what
// type of address to search for. The close to fields are normally treated as a
// receiver, if you would like to exclude them set this parameter to true.
func (s *LookupAssetTransactions) ExcludeCloseTo(ExcludeCloseTo bool) *LookupAssetTransactions {
	s.p.ExcludeCloseTo = ExcludeCloseTo
	return s
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *LookupAssetTransactions) Limit(Limit uint64) *LookupAssetTransactions {
	s.p.Limit = Limit
	return s
}

// MaxRound include results at or before the specified max-round.
func (s *LookupAssetTransactions) MaxRound(MaxRound uint64) *LookupAssetTransactions {
	s.p.MaxRound = MaxRound
	return s
}

// MinRound include results at or after the specified min-round.
func (s *LookupAssetTransactions) MinRound(MinRound uint64) *LookupAssetTransactions {
	s.p.MinRound = MinRound
	return s
}

// NextToken the next page of results. Use the next token provided by the previous
// results.
func (s *LookupAssetTransactions) NextToken(NextToken string) *LookupAssetTransactions {
	s.p.NextToken = NextToken
	return s
}

// NotePrefix specifies a prefix which must be contained in the note field.
func (s *LookupAssetTransactions) NotePrefix(NotePrefix []byte) *LookupAssetTransactions {
	s.p.NotePrefix = base64.StdEncoding.EncodeToString(NotePrefix)

	return s
}

// RekeyTo include results which include the rekey-to field.
func (s *LookupAssetTransactions) RekeyTo(RekeyTo bool) *LookupAssetTransactions {
	s.p.RekeyTo = RekeyTo
	return s
}

// Round include results for the specified round.
func (s *LookupAssetTransactions) Round(Round uint64) *LookupAssetTransactions {
	s.p.Round = Round
	return s
}

// SigType sigType filters just results using the specified type of signature:
// * sig - Standard
// * msig - MultiSig
// * lsig - LogicSig
func (s *LookupAssetTransactions) SigType(SigType string) *LookupAssetTransactions {
	s.p.SigType = SigType
	return s
}

// TxType
func (s *LookupAssetTransactions) TxType(TxType string) *LookupAssetTransactions {
	s.p.TxType = TxType
	return s
}

// TXID lookup the specific transaction by ID.
func (s *LookupAssetTransactions) TXID(TXID string) *LookupAssetTransactions {
	s.p.TXID = TXID
	return s
}

// Do performs the HTTP request
func (s *LookupAssetTransactions) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/assets/%v/transactions", s.assetId), s.p, headers)
	return
}
