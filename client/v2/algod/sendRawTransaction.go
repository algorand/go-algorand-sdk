package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"strings"
)

type SendRawTransaction struct {
	c   *Client
	stx []byte
}

type txidresponse struct {
	TxID string `json:"txId"`
}

func (s *SendRawTransaction) Do(ctx context.Context, headers ...*common.Header) (txid string, err error) {
	var response txidresponse
	// Set default Content-Type, if the user didn't specify it.
	addContentType := true
	for _, header := range headers {
		if strings.ToLower(header.Key) == "content-type" {
			addContentType = false
			break
		}
	}
	if addContentType {
		headers = append(headers, &common.Header{"Content-Type", "application/x-binary"})
	}
	err = s.c.post(ctx, &response, "/transactions", s.stx, headers)
	txid = response.TxID
	return
}
