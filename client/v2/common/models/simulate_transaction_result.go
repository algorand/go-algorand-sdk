package models

// SimulateTransactionResult simulation result for an individual transaction
type SimulateTransactionResult struct {
	// AppBudgetConsumed budget used during execution of an app call transaction. This
	// value includes budged used by inner app calls spawned by this transaction.
	AppBudgetConsumed uint64 `json:"app-budget-consumed,omitempty"`

	// ExecTrace the execution trace of calling an app or a logic sig, containing the
	// inner app call trace in a recursive way.
	ExecTrace SimulationTransactionExecTrace `json:"exec-trace,omitempty"`

	// LogicSigBudgetConsumed budget used during execution of a logic sig transaction.
	LogicSigBudgetConsumed uint64 `json:"logic-sig-budget-consumed,omitempty"`

	// TxnResult details about a pending transaction. If the transaction was recently
	// confirmed, includes confirmation details like the round and reward details.
	TxnResult PendingTransactionResponse `json:"txn-result"`

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
