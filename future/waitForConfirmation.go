package future

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// `WaitForConfirmation` waits for a pending transaction to be accepted by the network
// `txid`: The ID of the pending transaction to wait for
// `waitRounds`: The number of rounds to block before existing with an error. If zero, there is no timeout
func WaitForConfirmation(c *algod.Client, txid string, waitRounds uint64, ctx context.Context, headers ...*common.Header) (txInfo models.PendingTransactionInfoResponse, err error) {
	response, err := c.Status().Do(ctx, headers...)
	if err != nil {
		return
	}

	lastRound := response.LastRound
	currentRound := lastRound + 1

	for {
		// Check that the `waitRounds` has not passed
		if waitRounds > 0 && currentRound > lastRound+waitRounds {
			err = fmt.Errorf("Wait for transaction id %s timed out", txid)
			return
		}

		txInfo, _, err = c.PendingTransactionInformation(txid).Do(ctx, headers...)
		if err != nil {
			return
		}

		// The transaction has been confirmed
		if txInfo.ConfirmedRound > 0 {
			return
		}

		// Wait until the block for the `currentRound` is confirmed
		response, err = c.StatusAfterBlock(currentRound).Do(ctx, headers...)
		if err != nil {
			return
		}

		// Increment the `currentRound`
		currentRound += 1
	}

	return
}
