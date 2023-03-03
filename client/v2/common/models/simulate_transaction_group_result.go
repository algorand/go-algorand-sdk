package models

// SimulateTransactionGroupResult simulation result for an atomic transaction group
type SimulateTransactionGroupResult struct {
	// FailedAt if present, indicates which transaction in this group caused the
	// failure. This array represents the path to the failing transaction. Indexes are
	// zero based, the first element indicates the top-level transaction, and
	// successive elements indicate deeper inner transactions.
	FailedAt []uint64 `json:"failed-at,omitempty"`

	// FailureMessage if present, indicates that the transaction group failed and
	// specifies why that happened
	FailureMessage string `json:"failure-message,omitempty"`

	// TxnResults simulation result for individual transactions
	TxnResults []SimulateTransactionResult `json:"txn-results"`
}