package models

// SimulateRequest request type for simulation endpoint.
type SimulateRequest struct {
	// TxnGroups the transaction groups to simulate.
	TxnGroups []SimulateRequestTransactionGroup `json:"txn-groups"`
}
