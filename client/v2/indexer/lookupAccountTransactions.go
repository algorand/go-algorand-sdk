package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountTransactionsParams contains all of the query parameters for url serialization.
type LookupAccountTransactionsParams struct {

	// AfterTime include results after the given time. Must be an RFC 3339 formatted
	// string.
	AfterTime string `url:"after-time,omitempty"`

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

// LookupAccountTransactions lookup account transactions. Transactions are returned
// newest to oldest.
type LookupAccountTransactions struct {
	c *Client

	accountId string

	p LookupAccountTransactionsParams
}

// AfterTimeString include results after the given time. Must be an RFC 3339
// formatted string.
func (s *LookupAccountTransactions) AfterTimeString(AfterTime string) *LookupAccountTransactions {
	s.p.AfterTime = AfterTime
	return s
}

// AfterTime include results after the given time. Must be an RFC 3339 formatted
// string.
func (s *LookupAccountTransactions) AfterTime(AfterTime time.Time) *LookupAccountTransactions {
	AfterTimeStr := AfterTime.Format(time.RFC3339)

	return s.AfterTimeString(AfterTimeStr)
}

// AssetID asset ID
func (s *LookupAccountTransactions) AssetID(AssetID uint64) *LookupAccountTransactions {
	s.p.AssetID = AssetID
	return s
}

// BeforeTimeString include results before the given time. Must be an RFC 3339
// formatted string.
func (s *LookupAccountTransactions) BeforeTimeString(BeforeTime string) *LookupAccountTransactions {
	s.p.BeforeTime = BeforeTime
	return s
}

// BeforeTime include results before the given time. Must be an RFC 3339 formatted
// string.
func (s *LookupAccountTransactions) BeforeTime(BeforeTime time.Time) *LookupAccountTransactions {
	BeforeTimeStr := BeforeTime.Format(time.RFC3339)

	return s.BeforeTimeString(BeforeTimeStr)
}

// CurrencyGreaterThan results should have an amount greater than this value.
// MicroAlgos are the default currency unless an asset-id is provided, in which
// case the asset will be used.
func (s *LookupAccountTransactions) CurrencyGreaterThan(CurrencyGreaterThan uint64) *LookupAccountTransactions {
	s.p.CurrencyGreaterThan = CurrencyGreaterThan
	return s
}

// CurrencyLessThan results should have an amount less than this value. MicroAlgos
// are the default currency unless an asset-id is provided, in which case the asset
// will be used.
func (s *LookupAccountTransactions) CurrencyLessThan(CurrencyLessThan uint64) *LookupAccountTransactions {
	s.p.CurrencyLessThan = CurrencyLessThan
	return s
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *LookupAccountTransactions) Limit(Limit uint64) *LookupAccountTransactions {
	s.p.Limit = Limit
	return s
}

// MaxRound include results at or before the specified max-round.
func (s *LookupAccountTransactions) MaxRound(MaxRound uint64) *LookupAccountTransactions {
	s.p.MaxRound = MaxRound
	return s
}

// MinRound include results at or after the specified min-round.
func (s *LookupAccountTransactions) MinRound(MinRound uint64) *LookupAccountTransactions {
	s.p.MinRound = MinRound
	return s
}

// NextToken the next page of results. Use the next token provided by the previous
// results.
func (s *LookupAccountTransactions) NextToken(NextToken string) *LookupAccountTransactions {
	s.p.NextToken = NextToken
	return s
}

// NotePrefix specifies a prefix which must be contained in the note field.
func (s *LookupAccountTransactions) NotePrefix(NotePrefix []byte) *LookupAccountTransactions {
	s.p.NotePrefix = base64.StdEncoding.EncodeToString(NotePrefix)

	return s
}

// RekeyTo include results which include the rekey-to field.
func (s *LookupAccountTransactions) RekeyTo(RekeyTo bool) *LookupAccountTransactions {
	s.p.RekeyTo = RekeyTo
	return s
}

// Round include results for the specified round.
func (s *LookupAccountTransactions) Round(Round uint64) *LookupAccountTransactions {
	s.p.Round = Round
	return s
}

// SigType sigType filters just results using the specified type of signature:
// * sig - Standard
// * msig - MultiSig
// * lsig - LogicSig
func (s *LookupAccountTransactions) SigType(SigType string) *LookupAccountTransactions {
	s.p.SigType = SigType
	return s
}

// TxType
func (s *LookupAccountTransactions) TxType(TxType string) *LookupAccountTransactions {
	s.p.TxType = TxType
	return s
}

// TXID lookup the specific transaction by ID.
func (s *LookupAccountTransactions) TXID(TXID string) *LookupAccountTransactions {
	s.p.TXID = TXID
	return s
}

// Do performs the HTTP request
func (s *LookupAccountTransactions) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%v/transactions", s.accountId), s.p, headers)
	return
}
