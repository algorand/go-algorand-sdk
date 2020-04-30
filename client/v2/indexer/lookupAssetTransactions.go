package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type LookupAssetTransactions struct {
	c     *Client
	index uint64
	p     models.LookupAssetTransactionsParams
}

func (s *LookupAssetTransactions) NextToken(nextToken string) *LookupAssetTransactions {
	s.p.NextToken = nextToken
	return s
}

func (s *LookupAssetTransactions) NotePrefix(prefix []byte) *LookupAssetTransactions {
	s.p.NotePrefix = base64.StdEncoding.EncodeToString(prefix)
	return s
}

func (s *LookupAssetTransactions) TxType(txtype string) *LookupAssetTransactions {
	s.p.TxType = txtype
	return s
}

func (s *LookupAssetTransactions) SigType(sigtype string) *LookupAssetTransactions {
	s.p.SigType = sigtype
	return s
}

func (s *LookupAssetTransactions) TXID(txid string) *LookupAssetTransactions {
	s.p.TxId = txid
	return s
}

func (s *LookupAssetTransactions) Round(rnd uint64) *LookupAssetTransactions {
	s.p.Round = rnd
	return s
}

func (s *LookupAssetTransactions) MinRound(min uint64) *LookupAssetTransactions {
	s.p.MinRound = min
	return s
}

func (s *LookupAssetTransactions) MaxRound(max uint64) *LookupAssetTransactions {
	s.p.MaxRound = max
	return s
}

func (s *LookupAssetTransactions) Address(address string) *LookupAssetTransactions {
	s.p.Address = address
	return s
}

func (s *LookupAssetTransactions) Limit(limit uint64) *LookupAssetTransactions {
	s.p.Limit = limit
	return s
}

func (s *LookupAssetTransactions) BeforeTime(before string) *LookupAssetTransactions {
	s.p.BeforeTime = before
	return s
}

func (s *LookupAssetTransactions) AfterTime(after string) *LookupAssetTransactions {
	s.p.AfterTime = after
	return s
}

func (s *LookupAssetTransactions) CurrencyGreaterThan(greaterThan uint64) *LookupAssetTransactions {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *LookupAssetTransactions) CurrencyLessThan(lessThan uint64) *LookupAssetTransactions {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *LookupAssetTransactions) AddressRole(role string) *LookupAssetTransactions {
	s.p.AddressRole = role
	return s
}

func (s *LookupAssetTransactions) ExcludeCloseTo(exclude bool) *LookupAssetTransactions {
	s.p.ExcludeCloseTo = exclude
	return s
}

func (s *LookupAssetTransactions) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, transactions []models.Transaction, err error) {
	var response models.TransactionsResponse
	err = s.c.get(ctx, &response, fmt.Sprintf("/assets/%d/transactions", s.index), s.p, headers)
	validRound = response.CurrentRound
	transactions = response.Transactions
	return
}
