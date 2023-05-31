package transaction

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// WaitForConfirmation waits for a pending transaction to be accepted by the network
// txid: The ID of the pending transaction to wait for
// waitRounds: The number of rounds to block before exiting with an error.
func WaitForConfirmation(c *algod.Client, txid string, waitRounds uint64, ctx context.Context, headers ...*common.Header) (txInfo models.PendingTransactionInfoResponse, err error) { //nolint:revive // Ignore Context order for backwards compatibility
	response, err := c.Status().Do(ctx, headers...)
	if err != nil {
		return
	}

	lastRound := response.LastRound
	currentRound := lastRound + 1

	for {
		// Check that the `waitRounds` has not passed
		if currentRound > lastRound+waitRounds {
			err = fmt.Errorf("Wait for transaction id %s timed out", txid)
			return
		}

		txInfo, _, err = c.PendingTransactionInformation(txid).Do(ctx, headers...)
		if err == nil {
			if len(txInfo.PoolError) != 0 {
				// The transaction has been rejected
				err = fmt.Errorf("Transaction rejected: %s", txInfo.PoolError)
				return
			}

			if txInfo.ConfirmedRound > 0 {
				// The transaction has been confirmed
				return
			}
		}
		// ignore errors from PendingTransactionInformation, since it may return 404 if the algod
		// instance is behind a load balancer and the request goes to a different algod than the
		// one we submitted the transaction to
		err = nil

		// Wait until the block for the `currentRound` is confirmed
		response, err = c.StatusAfterBlock(currentRound).Do(ctx, headers...)
		if err != nil {
			return
		}

		// Increment the `currentRound`
		currentRound++
	}
}
