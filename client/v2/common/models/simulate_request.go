package models

// SimulateRequest request type for simulation endpoint.
type SimulateRequest struct {
	// LiftLogLimits the boolean flag that lifts the limit on log opcode during
	// simulation.
	LiftLogLimits bool `json:"lift-log-limits,omitempty"`

	// TxnGroups the transaction groups to simulate.
	TxnGroups []SimulateRequestTransactionGroup `json:"txn-groups"`
}
