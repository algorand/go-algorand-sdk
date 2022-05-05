package models

// DryrunTxnResult dryrunTxnResult contains any LogicSig or ApplicationCall program
// debug information and state updates from a dryrun.
type DryrunTxnResult struct {
	// AppCallMessages
	AppCallMessages []string `json:"app-call-messages,omitempty"`

	// AppCallTrace
	AppCallTrace []DryrunState `json:"app-call-trace,omitempty"`

	// Cost execution cost of app call transaction
	Cost int64 `json:"cost,omitempty"`

	// Disassembly disassembled program line by line.
	Disassembly []string `json:"disassembly"`

	// GlobalDelta application state delta.
	GlobalDelta []EvalDeltaKeyValue `json:"global-delta,omitempty"`

	// LocalDeltas
	LocalDeltas []AccountStateDelta `json:"local-deltas,omitempty"`

	// LogicSigDisassembly disassembled lsig program line by line.
	LogicSigDisassembly []string `json:"logic-sig-disassembly,omitempty"`

	// LogicSigMessages
	LogicSigMessages []string `json:"logic-sig-messages,omitempty"`

	// LogicSigTrace
	LogicSigTrace []DryrunState `json:"logic-sig-trace,omitempty"`

	// Logs
	Logs [][]byte `json:"logs,omitempty"`
}
