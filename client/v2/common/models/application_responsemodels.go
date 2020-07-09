package models

import "github.com/algorand/go-algorand-sdk/types"

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

// ApplicationStateSchema specifies maximums on the number of each type that may be
// stored.
type ApplicationStateSchema struct {
	// NumByteSlice (nbs) num of byte slices.
	NumByteSlice uint64 `json:"num-byte-slice,omitempty"`

	// NumUint (nui) num of uints.
	NumUint uint64 `json:"num-uint,omitempty"`
}

// ApplicationLocalStates pair of application index and application local state
type ApplicationLocalStates struct {
	Id uint64 `json:"id,omitempty"`

	// State stores local state associated with an application.
	State ApplicationLocalState `json:"state,omitempty"`
}

// ApplicationLocalState stores local state associated with an application.
type ApplicationLocalState struct {
	// KeyValue (tkv) storage.
	KeyValue []TealKeyValue `json:"key-value,omitempty"`

	// Schema (hsch) schema.
	Schema ApplicationStateSchema `json:"schema,omitempty"`
}

// TealKeyValue represents a key-value pair in an application store.
type TealKeyValue struct {
	Key string `json:"key,omitempty"`

	// Value represents a TEAL value.
	Value TealValue `json:"value,omitempty"`
}

// TealValue represents a TEAL value.
type TealValue struct {
	// Bytes (tb) bytes value.
	Bytes string `json:"bytes,omitempty"`

	// Type (tt) value type.
	Type uint64 `json:"type,omitempty"`

	// Uint (ui) uint value.
	Uint uint64 `json:"uint,omitempty"`
}

// AccountStateDelta application state delta.
type AccountStateDelta struct {
	Address string `json:"address,omitempty"`

	// Delta application state delta.
	Delta []EvalDeltaKeyValue `json:"delta,omitempty"`
}

// EvalDeltaKeyValue key-value pairs for StateDelta.
type EvalDeltaKeyValue struct {
	Key string `json:"key,omitempty"`

	// Value represents a TEAL value delta.
	Value EvalDelta `json:"value,omitempty"`
}

// EvalDelta represents a TEAL value delta.
type EvalDelta struct {
	// Action (at) delta action.
	Action uint64 `json:"action,omitempty"`

	// Bytes (bs) bytes value.
	Bytes string `json:"bytes,omitempty"`

	// Uint (ui) uint value.
	Uint uint64 `json:"uint,omitempty"`
}

// Application index and its parameters
type Application struct {
	// Id (appidx) application index.
	Id uint64 `json:"id,omitempty"`

	// Params (appparams) application parameters.
	Params ApplicationParams `json:"params,omitempty"`
}

// ApplicationParams stores the global information associated with an application.
type ApplicationParams struct {
	// ApprovalProgram (approv) approval program.
	ApprovalProgram string `json:"approval-program,omitempty"`

	// ClearStateProgram (clearp) approval program.
	ClearStateProgram string `json:"clear-state-program,omitempty"`

	// Creator the address that created this application. This is the address where the
	// parameters and global state for this application can be found.
	Creator string `json:"creator,omitempty"`

	// GlobalState [\gs) global schema
	GlobalState []TealKeyValue `json:"global-state,omitempty"`

	// GlobalStateSchema [\lsch) global schema
	GlobalStateSchema ApplicationStateSchema `json:"global-state-schema,omitempty"`

	// LocalStateSchema [\lsch) local schema
	LocalStateSchema ApplicationStateSchema `json:"local-state-schema,omitempty"`
}

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

// CompileResponse is response from /v2/compile
type CompileResponse struct {
	// Hash base32 SHA512_256 of program bytes (Address style)
	Hash string `json:"hash,omitempty"`

	// Result base64 encoded program bytes
	Result string `json:"result,omitempty"`
}

// DryrunResponse is response from /v2/dryrun
type DryrunResponse struct {
	Error string `json:"error,omitempty"`

	// ProtocolVersion protocol version is the protocol version Dryrun was operated
	// under.
	ProtocolVersion string `json:"protocol-version,omitempty"`

	Txns []DryrunTxnResult `json:"txns,omitempty"`
}
