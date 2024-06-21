package models

// SimulateRequest request type for simulation endpoint.
type SimulateRequest struct {
	// AllowEmptySignatures allows transactions without signatures to be simulated as
	// if they had correct signatures.
	AllowEmptySignatures bool `json:"allow-empty-signatures,omitempty"`

	// AllowMoreLogging lifts limits on log opcode usage during simulation.
	AllowMoreLogging bool `json:"allow-more-logging,omitempty"`

	// AllowUnnamedResources allows access to unnamed resources during simulation.
	AllowUnnamedResources bool `json:"allow-unnamed-resources,omitempty"`

	// ExecTraceConfig an object that configures simulation execution trace.
	ExecTraceConfig SimulateTraceConfig `json:"exec-trace-config,omitempty"`

	// ExtraOpcodeBudget applies extra opcode budget during simulation for each
	// transaction group.
	ExtraOpcodeBudget uint64 `json:"extra-opcode-budget,omitempty"`

	// FixSigners if true, signers for transactions that are missing signatures will be
	// fixed during evaluation.
	FixSigners bool `json:"fix-signers,omitempty"`

	// Round if provided, specifies the round preceding the simulation. State changes
	// through this round will be used to run this simulation. Usually only the 4 most
	// recent rounds will be available (controlled by the node config value
	// MaxAcctLookback). If not specified, defaults to the latest available round.
	Round uint64 `json:"round,omitempty"`

	// TxnGroups the transaction groups to simulate.
	TxnGroups []SimulateRequestTransactionGroup `json:"txn-groups"`
}
