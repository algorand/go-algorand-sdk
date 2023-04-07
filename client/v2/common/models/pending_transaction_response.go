package models

import "github.com/algorand/go-algorand-sdk/v2/types"

// PendingTransactionResponse details about a pending transaction. If the
// transaction was recently confirmed, includes confirmation details like the round
// and reward details.
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

	// GlobalStateDelta global state key/value changes for the application being
	// executed by this transaction.
	GlobalStateDelta []EvalDeltaKeyValue `json:"global-state-delta,omitempty"`

	// InnerTxns inner transactions produced by application execution.
	InnerTxns []PendingTransactionResponse `json:"inner-txns,omitempty"`

	// LocalStateDelta local state key/value changes for the application being executed
	// by this transaction.
	LocalStateDelta []AccountStateDelta `json:"local-state-delta,omitempty"`

	// Logs logs for the application being executed by this transaction.
	Logs [][]byte `json:"logs,omitempty"`

	// PoolError indicates that the transaction was kicked out of this node's
	// transaction pool (and specifies why that happened). An empty string indicates
	// the transaction wasn't kicked out of this node's txpool due to an error.
	PoolError string `json:"pool-error"`

	// ReceiverRewards rewards in microalgos applied to the receiver account.
	ReceiverRewards uint64 `json:"receiver-rewards,omitempty"`

	// SenderRewards rewards in microalgos applied to the sender account.
	SenderRewards uint64 `json:"sender-rewards,omitempty"`

	// Transaction the raw signed transaction.
	Transaction types.SignedTxn `json:"txn"`
}
