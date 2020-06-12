package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

/**
 * /v2/assets/{asset-id}/transactions
 * Lookup transactions for an asset.
 */
type LookupAssetTransactions struct {
	c       *Client
	p       models.LookupAssetTransactionsParams
	assetId uint64
}

/**
 * Only include transactions with this address in one of the transaction fields.
 */
func (s *LookupAssetTransactions) Address(address types.Address) *LookupAssetTransactions {
	s.p.Address = address.String()
	return s
}

func (s *LookupAssetTransactions) AddressString(address string) *LookupAssetTransactions {
	s.p.Address = address
	return s
}

/**
 * Combine with the address parameter to define what type of address to search for.
 */
func (s *LookupAssetTransactions) AddressRole(addressRole string) *LookupAssetTransactions {
	s.p.AddressRole = addressRole
	return s
}

/**
 * Include results after the given time. Must be an RFC 3339 formatted string.
 */
func (s *LookupAssetTransactions) AfterTime(afterTime time.Time) *LookupAssetTransactions {
	s.p.AfterTime = afterTime.Format(time.RFC3339)
	return s
}

func (s *LookupAssetTransactions) AfterTimeString(after string) *LookupAssetTransactions {
	s.p.AfterTime = after
	return s
}

/**
 * Application ID
 */
func (s *LookupAssetTransactions) ApplicationId(applicationId uint64) *LookupAssetTransactions {
	s.p.ApplicationId = applicationId
	return s
}

/**
 * Include results before the given time. Must be an RFC 3339 formatted string.
 */
func (s *LookupAssetTransactions) BeforeTime(beforeTime time.Time) *LookupAssetTransactions {
	s.p.BeforeTime = beforeTime.Format(time.RFC3339)
	return s
}

func (s *LookupAssetTransactions) BeforeTimeString(before string) *LookupAssetTransactions {
	s.p.BeforeTime = before
	return s
}

/**
 * Results should have an amount greater than this value. MicroAlgos are the
 * default currency unless an asset-id is provided, in which case the asset will be
 * used.
 */
func (s *LookupAssetTransactions) CurrencyGreaterThan(currencyGreaterThan uint64) *LookupAssetTransactions {
	s.p.CurrencyGreaterThan = currencyGreaterThan
	return s
}

/**
 * Results should have an amount less than this value. MicroAlgos are the default
 * currency unless an asset-id is provided, in which case the asset will be used.
 */
func (s *LookupAssetTransactions) CurrencyLessThan(currencyLessThan uint64) *LookupAssetTransactions {
	s.p.CurrencyLessThan = currencyLessThan
	return s
}

/**
 * Combine with address and address-role parameters to define what type of address
 * to search for. The close to fields are normally treated as a receiver, if you
 * would like to exclude them set this parameter to true.
 */
func (s *LookupAssetTransactions) ExcludeCloseTo(excludeCloseTo bool) *LookupAssetTransactions {
	s.p.ExcludeCloseTo = excludeCloseTo
	return s
}

/**
 * Maximum number of results to return.
 */
func (s *LookupAssetTransactions) Limit(limit uint64) *LookupAssetTransactions {
	s.p.Limit = limit
	return s
}

/**
 * Include results at or before the specified max-round.
 */
func (s *LookupAssetTransactions) MaxRound(maxRound uint64) *LookupAssetTransactions {
	s.p.MaxRound = maxRound
	return s
}

/**
 * Include results at or after the specified min-round.
 */
func (s *LookupAssetTransactions) MinRound(minRound uint64) *LookupAssetTransactions {
	s.p.MinRound = minRound
	return s
}

/**
 * The next page of results. Use the next token provided by the previous results.
 */
func (s *LookupAssetTransactions) Next(next string) *LookupAssetTransactions {
	s.p.Next = next
	return s
}

/**
 * Specifies a prefix which must be contained in the note field.
 */
func (s *LookupAssetTransactions) NotePrefix(notePrefix []byte) *LookupAssetTransactions {
	s.p.NotePrefix = base64.StdEncoding.EncodeToString(notePrefix)
	return s
}

/**
 * Include results which include the rekey-to field.
 */
func (s *LookupAssetTransactions) RekeyTo(rekeyTo bool) *LookupAssetTransactions {
	s.p.RekeyTo = rekeyTo
	return s
}

/**
 * Include results for the specified round.
 */
func (s *LookupAssetTransactions) Round(round uint64) *LookupAssetTransactions {
	s.p.Round = round
	return s
}

/**
 * SigType filters just results using the specified type of signature:
 *   sig - Standard
 *   msig - MultiSig
 *   lsig - LogicSig
 */
func (s *LookupAssetTransactions) SigType(sigType string) *LookupAssetTransactions {
	s.p.SigType = sigType
	return s
}

func (s *LookupAssetTransactions) TxType(txType string) *LookupAssetTransactions {
	s.p.TxType = txType
	return s
}

/**
 * Lookup the specific transaction by ID.
 */
func (s *LookupAssetTransactions) TXID(txid string) *LookupAssetTransactions {
	s.p.Txid = txid
	return s
}

func (s *LookupAssetTransactions) Do(ctx context.Context,
	headers ...*common.Header) (response models.TransactionsResponse, err error) {
	err = s.c.get(ctx, &response,
		fmt.Sprintf("/v2/assets/%d/transactions", s.assetId), s.p, headers)
	return
}
