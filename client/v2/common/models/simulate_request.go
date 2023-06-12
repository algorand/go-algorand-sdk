package models

// SimulateRequest request type for simulation endpoint.
type SimulateRequest struct {
	// AllowEmptySignatures allow transactions without signatures to be simulated as if
	// they had correct signatures.
	AllowEmptySignatures bool `json:"allow-empty-signatures,omitempty"`

	// AllowMoreLogging lifts limits on log opcode usage during simulation.
	AllowMoreLogging bool `json:"allow-more-logging,omitempty"`

	// ExecTraceConfig an object that configures simulation execution trace.
	ExecTraceConfig SimulateTraceConfig `json:"exec-trace-config,omitempty"`

	// ExtraOpcodeBudget applies extra opcode budget during simulation for each
	// transaction group.
	ExtraOpcodeBudget uint64 `json:"extra-opcode-budget,omitempty"`

	// TxnGroups the transaction groups to simulate.
	TxnGroups []SimulateRequestTransactionGroup `json:"txn-groups"`
}
