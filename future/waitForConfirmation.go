package future

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// WaitForConfirmation waits for a pending transaction to be accepted by the network
type WaitForConfirmation struct {
	c *algod.Client

	// The ID of the pending transaction to wait for
	txid string

	// The number of rounds to block before existing with an error. If zero, there is no timeout
	waitRounds uint64
}

// Do performs the HTTP request
func (s *WaitForConfirmation) Do(ctx context.Context, headers ...*common.Header) (response models.NodeStatus, err error) {
	response, err = s.c.Status().Do(ctx, headers...)
	if err != nil {
		return
	}

	lastRound := response.LastRound
	currentRound := lastRound + 1

	for {
		// Check that the `waitRounds` has not passed
		if s.waitRounds > 0 && currentRound > lastRound+s.waitRounds {
			err = fmt.Errorf("Wait for transaction id %s timed out", s.txid)
			return
		}

		var txInfo models.PendingTransactionInfoResponse
		txInfo, _, err = s.c.PendingTransactionInformation(s.txid).Do(ctx, headers...)
		if err != nil {
			return
		}

		// The transaction has been confirmed
		if txInfo.ConfirmedRound > 0 {
			return
		}

		// Wait until the block for the `currentRound` is confirmed
		response, err = s.c.StatusAfterBlock(currentRound).Do(ctx, headers...)
		if err != nil {
			return
		}

		// Increment the `currentRound`
		currentRound += 1
	}

	return
}
