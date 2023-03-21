package models

// SimulateTransactionResult simulation result for an individual transaction
type SimulateTransactionResult struct {
	// MissingSignature a boolean indicating whether this transaction is missing
	// signatures
	MissingSignature bool `json:"missing-signature,omitempty"`

	// TxnResult details about a pending transaction. If the transaction was recently
	// confirmed, includes confirmation details like the round and reward details.
	TxnResult PendingTransactionResponse `json:"txn-result"`
}
