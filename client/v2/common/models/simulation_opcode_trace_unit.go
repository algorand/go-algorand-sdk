package models

// SimulationOpcodeTraceUnit the set of trace information and effect from
// evaluating a single opcode.
type SimulationOpcodeTraceUnit struct {
	// Pc the program counter of the current opcode being evaluated.
	Pc uint64 `json:"pc"`

	// SpawnedInners the indexes of the traces for inner transactions spawned by this
	// opcode, if any.
	SpawnedInners []uint64 `json:"spawned-inners,omitempty"`
}
