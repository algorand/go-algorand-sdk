package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"strings"
)

type SendRawTransactionService struct {
	c   *Client
	stx []byte
}

func (s *SendRawTransactionService) Do(ctx context.Context, headers ...*common.Header) (txid string, err error) {
	var response models.TxId
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
	txid = string(response)
	return
}
