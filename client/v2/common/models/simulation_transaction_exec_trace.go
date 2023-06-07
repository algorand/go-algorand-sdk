package models

// SimulationTransactionExecTrace the execution trace of calling an app or a logic
// sig, containing the inner app call trace in a recursive way.
type SimulationTransactionExecTrace struct {
	// ApprovalProgramTrace program trace that contains a trace of opcode effects in an
	// approval program.
	ApprovalProgramTrace []SimulationOpcodeTraceUnit `json:"approval-program-trace,omitempty"`

	// ClearStateProgramTrace program trace that contains a trace of opcode effects in
	// a clear state program.
	ClearStateProgramTrace []SimulationOpcodeTraceUnit `json:"clear-state-program-trace,omitempty"`

	// InnerTrace an array of SimulationTransactionExecTrace representing the execution
	// trace of any inner transactions executed.
	InnerTrace []SimulationTransactionExecTrace `json:"inner-trace,omitempty"`

	// LogicSigTrace program trace that contains a trace of opcode effects in a logic
	// sig.
	LogicSigTrace []SimulationOpcodeTraceUnit `json:"logic-sig-trace,omitempty"`
}
