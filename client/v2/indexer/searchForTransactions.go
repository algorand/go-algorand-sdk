package indexer

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type SearchForTransactionsParams struct {

	// Address only include transactions with this address in one of the transaction
	// fields.
	Address string `url:"address,omitempty"`

	// AddressRole combine with the address parameter to define what type of address to
	// search for.
	AddressRole string `url:"address-role,omitempty"`

	// AfterTime include results after the given time. Must be an RFC 3339 formatted
	// string.
	AfterTime string `url:"after-time,omitempty"`

	// ApplicationId application ID
	ApplicationId uint64 `url:"application-id,omitempty"`

	// AssetID asset ID
	AssetID uint64 `url:"asset-id,omitempty"`

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

	// Limit maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// MaxRound include results at or before the specified max-round.
	MaxRound uint64 `url:"max-round,omitempty"`

	// MinRound include results at or after the specified min-round.
	MinRound uint64 `url:"min-round,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`

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

type SearchForTransactions struct {
	c *Client

	p SearchForTransactionsParams
}

// Address only include transactions with this address in one of the transaction
// fields.
func (s *SearchForTransactions) Address(Address string) *SearchForTransactions {
	s.p.Address = Address
	return s
}

// AddressRole combine with the address parameter to define what type of address to
// search for.
func (s *SearchForTransactions) AddressRole(AddressRole string) *SearchForTransactions {
	s.p.AddressRole = AddressRole
	return s
}

// AfterTimeString include results after the given time. Must be an RFC 3339
// formatted string.
func (s *SearchForTransactions) AfterTimeString(AfterTime string) *SearchForTransactions {
	s.p.AfterTime = AfterTime
	return s
}

// AfterTime include results after the given time. Must be an RFC 3339 formatted
// string.
func (s *SearchForTransactions) AfterTime(AfterTime time.Time) *SearchForTransactions {
	AfterTimeStr := AfterTime.Format(time.RFC3339)

	return s.AfterTimeString(AfterTimeStr)
}

// ApplicationId application ID
func (s *SearchForTransactions) ApplicationId(ApplicationId uint64) *SearchForTransactions {
	s.p.ApplicationId = ApplicationId
	return s
}

// AssetID asset ID
func (s *SearchForTransactions) AssetID(AssetID uint64) *SearchForTransactions {
	s.p.AssetID = AssetID
	return s
}

// BeforeTimeString include results before the given time. Must be an RFC 3339
// formatted string.
func (s *SearchForTransactions) BeforeTimeString(BeforeTime string) *SearchForTransactions {
	s.p.BeforeTime = BeforeTime
	return s
}

// BeforeTime include results before the given time. Must be an RFC 3339 formatted
// string.
func (s *SearchForTransactions) BeforeTime(BeforeTime time.Time) *SearchForTransactions {
	BeforeTimeStr := BeforeTime.Format(time.RFC3339)

	return s.BeforeTimeString(BeforeTimeStr)
}

// CurrencyGreaterThan results should have an amount greater than this value.
// MicroAlgos are the default currency unless an asset-id is provided, in which
// case the asset will be used.
func (s *SearchForTransactions) CurrencyGreaterThan(CurrencyGreaterThan uint64) *SearchForTransactions {
	s.p.CurrencyGreaterThan = CurrencyGreaterThan
	return s
}

// CurrencyLessThan results should have an amount less than this value. MicroAlgos
// are the default currency unless an asset-id is provided, in which case the asset
// will be used.
func (s *SearchForTransactions) CurrencyLessThan(CurrencyLessThan uint64) *SearchForTransactions {
	s.p.CurrencyLessThan = CurrencyLessThan
	return s
}

// ExcludeCloseTo combine with address and address-role parameters to define what
// type of address to search for. The close to fields are normally treated as a
// receiver, if you would like to exclude them set this parameter to true.
func (s *SearchForTransactions) ExcludeCloseTo(ExcludeCloseTo bool) *SearchForTransactions {
	s.p.ExcludeCloseTo = ExcludeCloseTo
	return s
}

// Limit maximum number of results to return.
func (s *SearchForTransactions) Limit(Limit uint64) *SearchForTransactions {
	s.p.Limit = Limit
	return s
}

// MaxRound include results at or before the specified max-round.
func (s *SearchForTransactions) MaxRound(MaxRound uint64) *SearchForTransactions {
	s.p.MaxRound = MaxRound
	return s
}

// MinRound include results at or after the specified min-round.
func (s *SearchForTransactions) MinRound(MinRound uint64) *SearchForTransactions {
	s.p.MinRound = MinRound
	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *SearchForTransactions) Next(Next string) *SearchForTransactions {
	s.p.Next = Next
	return s
}

// NotePrefix specifies a prefix which must be contained in the note field.
func (s *SearchForTransactions) NotePrefix(NotePrefix string) *SearchForTransactions {
	s.p.NotePrefix = NotePrefix
	return s
}

// RekeyTo include results which include the rekey-to field.
func (s *SearchForTransactions) RekeyTo(RekeyTo bool) *SearchForTransactions {
	s.p.RekeyTo = RekeyTo
	return s
}

// Round include results for the specified round.
func (s *SearchForTransactions) Round(Round uint64) *SearchForTransactions {
	s.p.Round = Round
	return s
}

// SigType sigType filters just results using the specified type of signature:
// * sig - Standard
// * msig - MultiSig
// * lsig - LogicSig
func (s *SearchForTransactions) SigType(SigType string) *SearchForTransactions {
	s.p.SigType = SigType
	return s
}

// TxType
func (s *SearchForTransactions) TxType(TxType string) *SearchForTransactions {
	s.p.TxType = TxType
	return s
}

// TXID lookup the specific transaction by ID.
func (s *SearchForTransactions) TXID(TXID string) *SearchForTransactions {
	s.p.TXID = TXID
	return s
}

func (s *SearchForTransactions) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/transactions", s.p, headers)
	return
}
