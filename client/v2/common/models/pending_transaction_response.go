package models

import "github.com/algorand/go-algorand-sdk/types"

// PendingTransactionResponse given a transaction id of a recently submitted
// transaction, it returns information about it. There are several cases when this
// might succeed:
// - transaction committed (committed round > 0)
// - transaction still in the pool (committed round = 0, pool error = "")
// - transaction removed from pool due to error (committed round = 0, pool error !=
// "")
// Or the transaction may have happened sufficiently long ago that the node no
// longer remembers it, and this will return an error.
type PendingTransactionResponse struct {
	// ApplicationIndex the application index if the transaction was found and it
	// created an application.
	ApplicationIndex uint64 `json:"application-index,omitempty"`

	// AssetClosingAmount the number of the asset's unit that were transferred to the
	// close-to address.
	AssetClosingAmount uint64 `json:"asset-closing-amount,omitempty"`

	// AssetIndex the asset index if the transaction was found and it created an asset.
	AssetIndex uint64 `json:"asset-index,omitempty"`

	// CloseRewards rewards in microalgos applied to the close remainder to account.
	CloseRewards uint64 `json:"close-rewards,omitempty"`

	// ClosingAmount closing amount for the transaction.
	ClosingAmount uint64 `json:"closing-amount,omitempty"`

	// ConfirmedRound the round where this transaction was confirmed, if present.
	ConfirmedRound uint64 `json:"confirmed-round,omitempty"`

	// GlobalStateDelta (gd) Global state key/value changes for the application being
	// executed by this transaction.
	GlobalStateDelta []EvalDeltaKeyValue `json:"global-state-delta,omitempty"`

	// LocalStateDelta (ld) Local state key/value changes for the application being
	// executed by this transaction.
	LocalStateDelta []AccountStateDelta `json:"local-state-delta,omitempty"`

	// PoolError indicates that the transaction was kicked out of this node's
	// transaction pool (and specifies why that happened). An empty string indicates
	// the transaction wasn't kicked out of this node's txpool due to an error.
	PoolError string `json:"pool-error,omitempty"`

	// ReceiverRewards rewards in microalgos applied to the receiver account.
	ReceiverRewards uint64 `json:"receiver-rewards,omitempty"`

	// SenderRewards rewards in microalgos applied to the sender account.
	SenderRewards uint64 `json:"sender-rewards,omitempty"`

	// Transaction the raw signed transaction.
	Transaction types.SignedTxn `json:"txn,omitempty"`
}
