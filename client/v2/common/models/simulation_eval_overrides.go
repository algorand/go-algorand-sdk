package models

// SimulationEvalOverrides the set of parameters and limits override during
// simulation. If this set of parameters is present, then evaluation parameters may
// differ from standard evaluation in certain ways.
type SimulationEvalOverrides struct {
	// MaxLogCalls the maximum log calls one can make during simulation
	MaxLogCalls uint64 `json:"max-log-calls,omitempty"`

	// MaxLogSize the maximum byte number to log during simulation
	MaxLogSize uint64 `json:"max-log-size,omitempty"`
}
