package models

// SimulateTransactionGroupResult simulation result for an atomic transaction group
type SimulateTransactionGroupResult struct {
	// AppBudgetAdded total budget added during execution of app calls in the
	// transaction group.
	AppBudgetAdded uint64 `json:"app-budget-added,omitempty"`

	// AppBudgetConsumed total budget consumed during execution of app calls in the
	// transaction group.
	AppBudgetConsumed uint64 `json:"app-budget-consumed,omitempty"`

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

	// UnnamedResourcesAccessed these are resources that were accessed by this group
	// that would normally have caused failure, but were allowed in simulation.
	// Depending on where this object is in the response, the unnamed resources it
	// contains may or may not qualify for group resource sharing. If this is a field
	// in SimulateTransactionGroupResult, the resources do qualify, but if this is a
	// field in SimulateTransactionResult, they do not qualify. In order to make this
	// group valid for actual submission, resources that qualify for group sharing can
	// be made available by any transaction of the group; otherwise, resources must be
	// placed in the same transaction which accessed them.
	UnnamedResourcesAccessed SimulateUnnamedResourcesAccessed `json:"unnamed-resources-accessed,omitempty"`
}
