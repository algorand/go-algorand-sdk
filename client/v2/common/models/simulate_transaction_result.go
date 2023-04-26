package models

// SimulateTransactionResult simulation result for an individual transaction
type SimulateTransactionResult struct {
	// AppBudgetConsumed budget used during execution of an app call transaction. This
	// value includes budged used by inner app calls spawned by this transaction.
	AppBudgetConsumed uint64 `json:"app-budget-consumed,omitempty"`

	// LogicSigBudgetConsumed budget used during execution of a logic sig transaction.
	LogicSigBudgetConsumed uint64 `json:"logic-sig-budget-consumed,omitempty"`

	// MissingSignature a boolean indicating whether this transaction is missing
	// signatures
	MissingSignature bool `json:"missing-signature,omitempty"`

	// TxnResult details about a pending transaction. If the transaction was recently
	// confirmed, includes confirmation details like the round and reward details.
	TxnResult PendingTransactionResponse `json:"txn-result"`
}
