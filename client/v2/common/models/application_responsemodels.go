package models

import "github.com/algorand/go-algorand-sdk/types"

// DryrunState stores the TEAL eval step data
type DryrunState struct {
	// Error evaluation error if any
	Error string `json:"error,omitempty"`

	// Line number
	Line uint64 `json:"line,omitempty"`

	// Pc program counter
	Pc uint64 `json:"pc,omitempty"`

	Scratch []TealValue `json:"scratch,omitempty"`

	Stack []TealValue `json:"stack,omitempty"`
}

// DryrunTxnResult contains any LogicSig or ApplicationCall program debug
// information and state updates from a dryrun.
type DryrunTxnResult struct {
	AppCallMessages []string `json:"app-call-messages,omitempty"`

	AppCallTrace []DryrunState `json:"app-call-trace,omitempty"`

	// Disassembly disassembled program line by line.
	Disassembly []string `json:"disassembly,omitempty"`

	// GlobalDelta application state delta.
	GlobalDelta []EvalDeltaKeyValue `json:"global-delta,omitempty"`

	LocalDeltas []AccountStateDelta `json:"local-deltas,omitempty"`

	LogicSigMessages []string `json:"logic-sig-messages,omitempty"`

	LogicSigTrace []DryrunState `json:"logic-sig-trace,omitempty"`
}

// DryrunRequest request data type for dryrun endpoint. Given the Transactions and
// simulated ledger state upload, run TEAL scripts and return debugging
// information.
type DryrunRequest struct {
	Accounts []Account `json:"accounts,omitempty"`

	Apps []Application `json:"apps,omitempty"`

	// LatestTimestamp is available to some TEAL scripts. Defaults to the latest
	// confirmed timestamp this algod is attached to.
	LatestTimestamp uint64 `json:"latest-timestamp,omitempty"`

	// ProtocolVersion specifies a specific version string to operate under, otherwise
	// whatever the current protocol of the network this algod is running in.
	ProtocolVersion string `json:"protocol-version,omitempty"`

	// Round is available to some TEAL scripts. Defaults to the current round on the
	// network this algod is attached to.
	Round uint64 `json:"round,omitempty"`

	Sources []DryrunSource `json:"sources,omitempty"`

	Txns []types.SignedTxn `json:"txns,omitempty"`
}

// DryrunSource is TEAL source text that gets uploaded, compiled, and inserted into
// transactions or application state.
type DryrunSource struct {
	AppIndex uint64 `json:"app-index,omitempty"`

	// FieldName is what kind of sources this is. If lsig then it goes into the
	// transactions[this.TxnIndex].LogicSig. If approv or clearp it goes into the
	// Approval Program or Clear State Program of application[this.AppIndex].
	FieldName string `json:"field-name,omitempty"`

	Source string `json:"source,omitempty"`

	TxnIndex uint64 `json:"txn-index,omitempty"`
}

type CatchpointStartResponse struct {
	// CatchupMessage catchup start response string
	CatchupMessage string `json:"catchup-message,omitempty"`
}

type CatchpointAbortResponse struct {
	// CatchupMessage catchup abort response string
	CatchupMessage string `json:"catchup-message,omitempty"`
}

// CompileResponse teal compile Result
type CompileResponse struct {
	// Hash base32 SHA512_256 of program bytes (Address style)
	Hash string `json:"hash,omitempty"`

	// Result base64 encoded program bytes
	Result string `json:"result,omitempty"`
}

// DryrunResponse contains per-txn debug information from a dryrun.
type DryrunResponse struct {
	Error string `json:"error,omitempty"`

	// ProtocolVersion protocol version is the protocol version Dryrun was operated
	// under.
	ProtocolVersion string `json:"protocol-version,omitempty"`

	Txns []DryrunTxnResult `json:"txns,omitempty"`
}
