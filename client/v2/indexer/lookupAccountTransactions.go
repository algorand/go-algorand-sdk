package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

/**
 * /v2/accounts/{account-id}/transactions
 * Lookup account transactions.
 */
type LookupAccountTransactions struct {
	c         *Client
	p         models.LookupAccountTransactionsParams
	accountId string
}

/**
 * Include results after the given time. Must be an RFC 3339 formatted string.
 */
func (s *LookupAccountTransactions) AfterTime(afterTime time.Time) *LookupAccountTransactions {
	s.p.AfterTime = afterTime.Format(time.RFC3339)
	return s
}

func (s *LookupAccountTransactions) AfterTimeString(after string) *LookupAccountTransactions {
	s.p.AfterTime = after
	return s
}

/**
 * Application ID
 */
func (s *LookupAccountTransactions) ApplicationId(applicationId uint64) *LookupAccountTransactions {
	s.p.ApplicationId = applicationId
	return s
}

/**
 * Asset ID
 */
func (s *LookupAccountTransactions) AssetID(assetId uint64) *LookupAccountTransactions {
	s.p.AssetId = assetId
	return s
}

/**
 * Include results before the given time. Must be an RFC 3339 formatted string.
 */
func (s *LookupAccountTransactions) BeforeTime(beforeTime time.Time) *LookupAccountTransactions {
	s.p.BeforeTime = beforeTime.Format(time.RFC3339)
	return s
}

func (s *LookupAccountTransactions) BeforeTimeString(before string) *LookupAccountTransactions {
	s.p.BeforeTime = before
	return s
}

/**
 * Results should have an amount greater than this value. MicroAlgos are the
 * default currency unless an asset-id is provided, in which case the asset will be
 * used.
 */
func (s *LookupAccountTransactions) CurrencyGreaterThan(currencyGreaterThan uint64) *LookupAccountTransactions {
	s.p.CurrencyGreaterThan = currencyGreaterThan
	return s
}

/**
 * Results should have an amount less than this value. MicroAlgos are the default
 * currency unless an asset-id is provided, in which case the asset will be used.
 */
func (s *LookupAccountTransactions) CurrencyLessThan(currencyLessThan uint64) *LookupAccountTransactions {
	s.p.CurrencyLessThan = currencyLessThan
	return s
}

/**
 * Maximum number of results to return.
 */
func (s *LookupAccountTransactions) Limit(limit uint64) *LookupAccountTransactions {
	s.p.Limit = limit
	return s
}

/**
 * Include results at or before the specified max-round.
 */
func (s *LookupAccountTransactions) MaxRound(maxRound uint64) *LookupAccountTransactions {
	s.p.MaxRound = maxRound
	return s
}

/**
 * Include results at or after the specified min-round.
 */
func (s *LookupAccountTransactions) MinRound(minRound uint64) *LookupAccountTransactions {
	s.p.MinRound = minRound
	return s
}

/**
 * The next page of results. Use the next token provided by the previous results.
 */
func (s *LookupAccountTransactions) Next(next string) *LookupAccountTransactions {
	s.p.Next = next
	return s
}

/**
 * Specifies a prefix which must be contained in the note field.
 */
func (s *LookupAccountTransactions) NotePrefix(notePrefix []byte) *LookupAccountTransactions {
	s.p.NotePrefix = base64.StdEncoding.EncodeToString(notePrefix)
	return s
}

/**
 * Include results which include the rekey-to field.
 */
func (s *LookupAccountTransactions) RekeyTo(rekeyTo bool) *LookupAccountTransactions {
	s.p.RekeyTo = rekeyTo
	return s
}

/**
 * Include results for the specified round.
 */
func (s *LookupAccountTransactions) Round(round uint64) *LookupAccountTransactions {
	s.p.Round = round
	return s
}

/**
 * SigType filters just results using the specified type of signature:
 *   sig - Standard
 *   msig - MultiSig
 *   lsig - LogicSig
 */
func (s *LookupAccountTransactions) SigType(sigType string) *LookupAccountTransactions {
	s.p.SigType = sigType
	return s
}

func (s *LookupAccountTransactions) TxType(txType string) *LookupAccountTransactions {
	s.p.TxType = txType
	return s
}

/**
 * Lookup the specific transaction by ID.
 */
func (s *LookupAccountTransactions) TXID(txid string) *LookupAccountTransactions {
	s.p.Txid = txid
	return s
}

func (s *LookupAccountTransactions) Do(ctx context.Context,
	headers ...*common.Header) (response models.TransactionsResponse, err error) {
	err = s.c.get(ctx, &response,
		fmt.Sprintf("/v2/accounts/%s/transactions", s.accountId), s.p, headers)
	return
}
