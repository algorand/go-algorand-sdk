package indexer

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"time"
)

type LookupAccountTransactionsService struct {
	c       *Client
	account string
	p       models.LookupAccountTransactionsParams
}

func (s *LookupAccountTransactionsService) NotePrefix(prefix []byte) *LookupAccountTransactionsService {
	s.p.NotePrefix = prefix
	return s
}

func (s *LookupAccountTransactionsService) TxType(txtype string) *LookupAccountTransactionsService {
	s.p.TxType = txtype
	return s
}

func (s *LookupAccountTransactionsService) SigType(sigtype string) *LookupAccountTransactionsService {
	s.p.SigType = sigtype
	return s
}

func (s *LookupAccountTransactionsService) TXID(txid string) *LookupAccountTransactionsService {
	s.p.TxId = txid
	return s
}

func (s *LookupAccountTransactionsService) Round(rnd uint64) *LookupAccountTransactionsService {
	s.p.Round = rnd
	return s
}

func (s *LookupAccountTransactionsService) MinRound(min uint64) *LookupAccountTransactionsService {
	s.p.MinRound = min
	return s
}

func (s *LookupAccountTransactionsService) MaxRound(max uint64) *LookupAccountTransactionsService {
	s.p.MaxRound = max
	return s
}

func (s *LookupAccountTransactionsService) AssetID(index uint64) *LookupAccountTransactionsService {
	s.p.AssetId = index
	return s
}

func (s *LookupAccountTransactionsService) Limit(limit uint64) *LookupAccountTransactionsService {
	s.p.Limit = limit
	return s
}

func (s *LookupAccountTransactionsService) BeforeTime(before time.Time) *LookupAccountTransactionsService {
	s.p.BeforeTime = before
	return s
}

func (s *LookupAccountTransactionsService) AfterTime(after time.Time) *LookupAccountTransactionsService {
	s.p.AfterTime = after
	return s
}

func (s *LookupAccountTransactionsService) CurrencyGreaterThan(greaterThan uint64) *LookupAccountTransactionsService {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *LookupAccountTransactionsService) CurrencyLessThan(lessThan uint64) *LookupAccountTransactionsService {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *LookupAccountTransactionsService) AddressRole(role string) *LookupAccountTransactionsService {
	s.p.AddressRole = role
	return s
}

func (s *LookupAccountTransactionsService) ExcludeCloseTo(exclude bool) *LookupAccountTransactionsService {
	s.p.ExcludeCloseTo = exclude
	return s
}

func (s *LookupAccountTransactionsService) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, transactions []models.Transaction, err error) {
	var response models.TransactionsResponse
	err = s.c.get(ctx, &response, fmt.Sprintf("/accounts/%s/transactions", s.account), s.p, headers)
	validRound = response.CurrentRound
	transactions = response.Transactions
	return
}
