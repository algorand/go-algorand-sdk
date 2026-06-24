package models

// SimulateResponse result of a transaction group simulation.
type SimulateResponse struct {
	// EvalOverrides the set of parameters and limits override during simulation. If
	// this set of parameters is present, then evaluation parameters may differ from
	// standard evaluation in certain ways.
	EvalOverrides SimulationEvalOverrides `json:"eval-overrides,omitempty"`

	// ExecTraceConfig an object that configures simulation execution trace.
	ExecTraceConfig SimulateTraceConfig `json:"exec-trace-config,omitempty"`

	// InitialStates initial states of resources that were accessed during simulation.
	InitialStates SimulateInitialStates `json:"initial-states,omitempty"`

	// LastRound the round immediately preceding this simulation. State changes through
	// this round were used to run this simulation.
	LastRound uint64 `json:"last-round"`

	// TotalFeesPaid total fees paid across all top-level transaction groups and their
	// descendants.
	TotalFeesPaid uint64 `json:"total-fees-paid,omitempty"`

	// TotalUsage total fee usage across all top-level transaction groups and their
	// descendants, in millionths of a basic transaction fee unit.
	TotalUsage uint64 `json:"total-usage,omitempty"`

	// TxnGroups a result object for each transaction group that was simulated.
	TxnGroups []SimulateTransactionGroupResult `json:"txn-groups"`

	// Version the version of this response object.
	Version uint64 `json:"version"`
}
