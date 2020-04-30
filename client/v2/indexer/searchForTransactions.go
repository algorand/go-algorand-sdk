package indexer

import (
	"context"
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type SearchForTransactions struct {
	c *Client
	p models.SearchForTransactionsParams
}

func (s *SearchForTransactions) NextToken(nextToken string) *SearchForTransactions {
	s.p.NextToken = nextToken
	return s
}

func (s *SearchForTransactions) NotePrefix(prefix []byte) *SearchForTransactions {
	s.p.NotePrefix = base64.StdEncoding.EncodeToString(prefix)
	return s
}

func (s *SearchForTransactions) TxType(txtype string) *SearchForTransactions {
	s.p.TxType = txtype
	return s
}

func (s *SearchForTransactions) SigType(sigtype string) *SearchForTransactions {
	s.p.SigType = sigtype
	return s
}

func (s *SearchForTransactions) TXID(txid string) *SearchForTransactions {
	s.p.TxId = txid
	return s
}

func (s *SearchForTransactions) Round(rnd uint64) *SearchForTransactions {
	s.p.Round = rnd
	return s
}

func (s *SearchForTransactions) MinRound(min uint64) *SearchForTransactions {
	s.p.MinRound = min
	return s
}

func (s *SearchForTransactions) MaxRound(max uint64) *SearchForTransactions {
	s.p.MaxRound = max
	return s
}

func (s *SearchForTransactions) AssetID(index uint64) *SearchForTransactions {
	s.p.AssetId = index
	return s
}

func (s *SearchForTransactions) Limit(limit uint64) *SearchForTransactions {
	s.p.Limit = limit
	return s
}

func (s *SearchForTransactions) BeforeTime(before string) *SearchForTransactions {
	s.p.BeforeTime = before
	return s
}

func (s *SearchForTransactions) AfterTime(after string) *SearchForTransactions {
	s.p.AfterTime = after
	return s
}

func (s *SearchForTransactions) CurrencyGreaterThan(greaterThan uint64) *SearchForTransactions {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *SearchForTransactions) CurrencyLessThan(lessThan uint64) *SearchForTransactions {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *SearchForTransactions) AddressRole(role string) *SearchForTransactions {
	s.p.AddressRole = role
	return s
}

func (s *SearchForTransactions) Address(address string) *SearchForTransactions {
	s.p.Address = address
	return s
}

func (s *SearchForTransactions) ExcludeCloseTo(exclude bool) *SearchForTransactions {
	s.p.ExcludeCloseTo = exclude
	return s
}

func (s *SearchForTransactions) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionsResponse, err error) {
	err = s.c.get(ctx, &response, "/transactions", s.p, headers)
	return
}
