package indexer

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"time"
)

type LookupAssetTransactionsService struct {
	c     *Client
	index uint64
	p     models.LookupAssetTransactionsParams
}

func (s *LookupAssetTransactionsService) NotePrefix(prefix []byte) *LookupAssetTransactionsService {
	s.p.NotePrefix = prefix
	return s
}

func (s *LookupAssetTransactionsService) TxType(txtype string) *LookupAssetTransactionsService {
	s.p.TxType = txtype
	return s
}

func (s *LookupAssetTransactionsService) SigType(sigtype string) *LookupAssetTransactionsService {
	s.p.SigType = sigtype
	return s
}

func (s *LookupAssetTransactionsService) TXID(txid string) *LookupAssetTransactionsService {
	s.p.TxId = txid
	return s
}

func (s *LookupAssetTransactionsService) Round(rnd uint64) *LookupAssetTransactionsService {
	s.p.Round = rnd
	return s
}

func (s *LookupAssetTransactionsService) MinRound(min uint64) *LookupAssetTransactionsService {
	s.p.MinRound = min
	return s
}

func (s *LookupAssetTransactionsService) MaxRound(max uint64) *LookupAssetTransactionsService {
	s.p.MaxRound = max
	return s
}

func (s *LookupAssetTransactionsService) Address(address string) *LookupAssetTransactionsService {
	s.p.Address = address
	return s
}

func (s *LookupAssetTransactionsService) Limit(limit uint64) *LookupAssetTransactionsService {
	s.p.Limit = limit
	return s
}

func (s *LookupAssetTransactionsService) BeforeTime(before time.Time) *LookupAssetTransactionsService {
	s.p.BeforeTime = before
	return s
}

func (s *LookupAssetTransactionsService) AfterTime(after time.Time) *LookupAssetTransactionsService {
	s.p.AfterTime = after
	return s
}

func (s *LookupAssetTransactionsService) CurrencyGreaterThan(greaterThan uint64) *LookupAssetTransactionsService {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *LookupAssetTransactionsService) CurrencyLessThan(lessThan uint64) *LookupAssetTransactionsService {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *LookupAssetTransactionsService) AddressRole(role string) *LookupAssetTransactionsService {
	s.p.AddressRole = role
	return s
}

func (s *LookupAssetTransactionsService) ExcludeCloseTo(exclude bool) *LookupAssetTransactionsService {
	s.p.ExcludeCloseTo = exclude
	return s
}

func (s *LookupAssetTransactionsService) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, transactions []models.Transaction, err error) {
	var response models.TransactionsResponse
	err = s.c.get(ctx, &response, fmt.Sprintf("/assets/%d/transactions", s.index), s.p, headers)
	validRound = response.CurrentRound
	transactions = response.Transactions
	return
}
