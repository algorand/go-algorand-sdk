package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"time"
)

type SearchForTransactionsService struct {
	c *Client
	p models.SearchForTransactionsParams
}

func (s *SearchForTransactionsService) NotePrefix(prefix []byte) *SearchForTransactionsService {
	s.p.NotePrefix = prefix
	return s
}

func (s *SearchForTransactionsService) TxType(txtype string) *SearchForTransactionsService {
	s.p.TxType = txtype
	return s
}

func (s *SearchForTransactionsService) SigType(sigtype string) *SearchForTransactionsService {
	s.p.SigType = sigtype
	return s
}

func (s *SearchForTransactionsService) TXID(txid string) *SearchForTransactionsService {
	s.p.TxId = txid
	return s
}

func (s *SearchForTransactionsService) Round(rnd uint64) *SearchForTransactionsService {
	s.p.Round = rnd
	return s
}

func (s *SearchForTransactionsService) MinRound(min uint64) *SearchForTransactionsService {
	s.p.MinRound = min
	return s
}

func (s *SearchForTransactionsService) MaxRound(max uint64) *SearchForTransactionsService {
	s.p.MaxRound = max
	return s
}

func (s *SearchForTransactionsService) AssetID(index uint64) *SearchForTransactionsService {
	s.p.AssetId = index
	return s
}

func (s *SearchForTransactionsService) Limit(limit uint64) *SearchForTransactionsService {
	s.p.Limit = limit
	return s
}

func (s *SearchForTransactionsService) BeforeTime(before time.Time) *SearchForTransactionsService {
	s.p.BeforeTime = before
	return s
}

func (s *SearchForTransactionsService) AfterTime(after time.Time) *SearchForTransactionsService {
	s.p.AfterTime = after
	return s
}

func (s *SearchForTransactionsService) CurrencyGreaterThan(greaterThan uint64) *SearchForTransactionsService {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *SearchForTransactionsService) CurrencyLessThan(lessThan uint64) *SearchForTransactionsService {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *SearchForTransactionsService) AddressRole(role string) *SearchForTransactionsService {
	s.p.AddressRole = role
	return s
}

func (s *SearchForTransactionsService) Address(address string) *SearchForTransactionsService {
	s.p.Address = address
	return s
}

func (s *SearchForTransactionsService) ExcludeCloseTo(exclude bool) *SearchForTransactionsService {
	s.p.ExcludeCloseTo = exclude
	return s
}

func (s *SearchForTransactionsService) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result []models.Transaction, err error) {
	var response models.TransactionsResponse
	err = s.c.get(ctx, &response, "/transactions", s.p, headers)
	validRound = response.CurrentRound
	result = response.Transactions
	return
}
