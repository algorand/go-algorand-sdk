package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

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

func (s *LookupAccountTransactions) BeforeTimeString(before string) *LookupAccountTransactions {
	s.p.BeforeTime = before
	return s
}

func (s *LookupAccountTransactions) AfterTimeString(after string) *LookupAccountTransactions {
	s.p.AfterTime = after
	return s
}

func (s *LookupAccountTransactions) BeforeTime(before time.Time) *LookupAccountTransactions {
	beforeString := before.Format(time.RFC3339)
	return s.BeforeTimeString(beforeString)
}

func (s *LookupAccountTransactions) AfterTime(after time.Time) *LookupAccountTransactions {
	afterString := after.Format(time.RFC3339)
	return s.AfterTimeString(afterString)
}

func (s *LookupAccountTransactions) CurrencyGreaterThan(greaterThan uint64) *LookupAccountTransactions {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *LookupAccountTransactions) CurrencyLessThan(lessThan uint64) *LookupAccountTransactions {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *LookupAccountTransactions) RekeyTo(rekeyTo bool) *LookupAccountTransactions {
	s.p.RekeyTo = rekeyTo
	return s
}

func (s *LookupAccountTransactions) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s/transactions", s.account), s.p, headers)
	return
}
