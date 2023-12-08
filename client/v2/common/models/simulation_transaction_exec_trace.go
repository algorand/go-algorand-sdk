package models

// SimulationTransactionExecTrace the execution trace of calling an app or a logic
// sig, containing the inner app call trace in a recursive way.
type SimulationTransactionExecTrace struct {
	// ApprovalProgramHash sHA512_256 hash digest of the approval program executed in
	// transaction.
	ApprovalProgramHash []byte `json:"approval-program-hash,omitempty"`

	// ApprovalProgramTrace program trace that contains a trace of opcode effects in an
	// approval program.
	ApprovalProgramTrace []SimulationOpcodeTraceUnit `json:"approval-program-trace,omitempty"`

	// ClearStateProgramHash sHA512_256 hash digest of the clear state program executed
	// in transaction.
	ClearStateProgramHash []byte `json:"clear-state-program-hash,omitempty"`

	// ClearStateProgramTrace program trace that contains a trace of opcode effects in
	// a clear state program.
	ClearStateProgramTrace []SimulationOpcodeTraceUnit `json:"clear-state-program-trace,omitempty"`

	// ClearStateRollback if true, indicates that the clear state program failed and
	// any persistent state changes it produced should be reverted once the program
	// exits.
	ClearStateRollback bool `json:"clear-state-rollback,omitempty"`

	// ClearStateRollbackError the error message explaining why the clear state program
	// failed. This field will only be populated if clear-state-rollback is true and
	// the failure was due to an execution error.
	ClearStateRollbackError string `json:"clear-state-rollback-error,omitempty"`

	// InnerTrace an array of SimulationTransactionExecTrace representing the execution
	// trace of any inner transactions executed.
	InnerTrace []SimulationTransactionExecTrace `json:"inner-trace,omitempty"`

	// LogicSigHash sHA512_256 hash digest of the logic sig executed in transaction.
	LogicSigHash []byte `json:"logic-sig-hash,omitempty"`

	// LogicSigTrace program trace that contains a trace of opcode effects in a logic
	// sig.
	LogicSigTrace []SimulationOpcodeTraceUnit `json:"logic-sig-trace,omitempty"`
}
