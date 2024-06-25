package models

// SimulationEvalOverrides the set of parameters and limits override during
// simulation. If this set of parameters is present, then evaluation parameters may
// differ from standard evaluation in certain ways.
type SimulationEvalOverrides struct {
	// AllowEmptySignatures if true, transactions without signatures are allowed and
	// simulated as if they were properly signed.
	AllowEmptySignatures bool `json:"allow-empty-signatures,omitempty"`

	// AllowUnnamedResources if true, allows access to unnamed resources during
	// simulation.
	AllowUnnamedResources bool `json:"allow-unnamed-resources,omitempty"`

	// ExtraOpcodeBudget the extra opcode budget added to each transaction group during
	// simulation
	ExtraOpcodeBudget uint64 `json:"extra-opcode-budget,omitempty"`

	// FixSigners if true, signers for transactions that are missing signatures will be
	// fixed during evaluation.
	FixSigners bool `json:"fix-signers,omitempty"`

	// MaxLogCalls the maximum log calls one can make during simulation
	MaxLogCalls uint64 `json:"max-log-calls,omitempty"`

	// MaxLogSize the maximum byte number to log during simulation
	MaxLogSize uint64 `json:"max-log-size,omitempty"`
}
