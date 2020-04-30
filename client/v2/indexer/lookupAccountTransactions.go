package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type LookupAccountTransactions struct {
	c       *Client
	account string
	p       models.LookupAccountTransactionsParams
}

func (s *LookupAccountTransactions) NextToken(nextToken string) *LookupAccountTransactions {
	s.p.NextToken = nextToken
	return s
}

func (s *LookupAccountTransactions) NotePrefix(prefix []byte) *LookupAccountTransactions {
	s.p.NotePrefix = base64.StdEncoding.EncodeToString(prefix)
	return s
}

func (s *LookupAccountTransactions) TxType(txtype string) *LookupAccountTransactions {
	s.p.TxType = txtype
	return s
}

func (s *LookupAccountTransactions) SigType(sigtype string) *LookupAccountTransactions {
	s.p.SigType = sigtype
	return s
}

func (s *LookupAccountTransactions) TXID(txid string) *LookupAccountTransactions {
	s.p.TxId = txid
	return s
}

func (s *LookupAccountTransactions) Round(rnd uint64) *LookupAccountTransactions {
	s.p.Round = rnd
	return s
}

func (s *LookupAccountTransactions) MinRound(min uint64) *LookupAccountTransactions {
	s.p.MinRound = min
	return s
}

func (s *LookupAccountTransactions) MaxRound(max uint64) *LookupAccountTransactions {
	s.p.MaxRound = max
	return s
}

func (s *LookupAccountTransactions) AssetID(index uint64) *LookupAccountTransactions {
	s.p.AssetId = index
	return s
}

func (s *LookupAccountTransactions) Limit(limit uint64) *LookupAccountTransactions {
	s.p.Limit = limit
	return s
}

func (s *LookupAccountTransactions) BeforeTime(before string) *LookupAccountTransactions {
	s.p.BeforeTime = before
	return s
}

func (s *LookupAccountTransactions) AfterTime(after string) *LookupAccountTransactions {
	s.p.AfterTime = after
	return s
}

func (s *LookupAccountTransactions) CurrencyGreaterThan(greaterThan uint64) *LookupAccountTransactions {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *LookupAccountTransactions) CurrencyLessThan(lessThan uint64) *LookupAccountTransactions {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *LookupAccountTransactions) AddressRole(role string) *LookupAccountTransactions {
	s.p.AddressRole = role
	return s
}

func (s *LookupAccountTransactions) ExcludeCloseTo(exclude bool) *LookupAccountTransactions {
	s.p.ExcludeCloseTo = exclude
	return s
}

func (s *LookupAccountTransactions) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, transactions []models.Transaction, err error) {
	var response models.TransactionsResponse
	err = s.c.get(ctx, &response, fmt.Sprintf("/accounts/%s/transactions", s.account), s.p, headers)
	validRound = response.CurrentRound
	transactions = response.Transactions
	return
}
