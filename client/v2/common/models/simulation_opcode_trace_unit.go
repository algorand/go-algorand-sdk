package models

// SimulationOpcodeTraceUnit the set of trace information and effect from
// evaluating a single opcode.
type SimulationOpcodeTraceUnit struct {
	// Pc the program counter of the current opcode being evaluated.
	Pc uint64 `json:"pc"`

	// SpawnedInners the indexes of the traces for inner transactions spawned by this
	// opcode, if any.
	SpawnedInners []uint64 `json:"spawned-inners,omitempty"`

	// StackAdditions the values added by this opcode to the stack.
	StackAdditions []AvmValue `json:"stack-additions,omitempty"`

	// StackPopCount the number of deleted stack values by this opcode.
	StackPopCount uint64 `json:"stack-pop-count,omitempty"`
}
