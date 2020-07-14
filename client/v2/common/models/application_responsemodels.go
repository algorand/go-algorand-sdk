package models

import "github.com/algorand/go-algorand-sdk/types"

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
	ApprovalProgram []byte `json:"approval-program,omitempty"`

	// ClearStateProgram (clearp) approval program.
	ClearStateProgram []byte `json:"clear-state-program,omitempty"`

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
// StateSchema represents a (apls) local-state or (apgs) global-state schema. These
// schemas determine how much storage may be used in a local-state or global-state
// for an application. The more space used, the larger minimum balance must be
// maintained in the account holding the data.
type StateSchema struct {
	// NumByteSlice maximum number of TEAL byte slices that may be stored in the
	// key/value store.
	NumByteSlice uint64 `json:"num-byte-slice,omitempty"`

	// NumUint maximum number of TEAL uints that may be stored in the key/value store.
	NumUint uint64 `json:"num-uint,omitempty"`
}

// TransactionApplication fields for application transactions.
// Definition:
// data/transactions/application.go : ApplicationCallTxnFields
type TransactionApplication struct {
	// Accounts (apat) List of accounts in addition to the sender that may be accessed
	// from the application's approval-program and clear-state-program.
	Accounts []string `json:"accounts,omitempty"`

	// ApplicationArgs (apaa) transaction specific arguments accessed from the
	// application's approval-program and clear-state-program.
	ApplicationArgs []byte `json:"application-args,omitempty"`

	// ApplicationId (apid) ID of the application being configured or empty if
	// creating.
	ApplicationId uint64 `json:"application-id,omitempty"`

	// ApprovalProgram (apap) Logic executed for every application transaction, except
	// when on-completion is set to "clear". It can read and write global state for the
	// application, as well as account-specific local state. Approval programs may
	// reject the transaction.
	ApprovalProgram string `json:"approval-program,omitempty"`

	// ClearStateProgram (apsu) Logic executed for application transactions with
	// on-completion set to "clear". It can read and write global state for the
	// application, as well as account-specific local state. Clear state programs
	// cannot reject the transaction.
	ClearStateProgram string `json:"clear-state-program,omitempty"`

	// ForeignApps (apfa) Lists the applications in addition to the application-id
	// whose global states may be accessed by this application's approval-program and
	// clear-state-program. The access is read-only.
	ForeignApps []uint64 `json:"foreign-apps,omitempty"`

	// GlobalStateSchema represents a (apls) local-state or (apgs) global-state schema.
	// These schemas determine how much storage may be used in a local-state or
	// global-state for an application. The more space used, the larger minimum balance
	// must be maintained in the account holding the data.
	GlobalStateSchema StateSchema `json:"global-state-schema,omitempty"`

	// LocalStateSchema represents a (apls) local-state or (apgs) global-state schema.
	// These schemas determine how much storage may be used in a local-state or
	// global-state for an application. The more space used, the larger minimum balance
	// must be maintained in the account holding the data.
	LocalStateSchema StateSchema `json:"local-state-schema,omitempty"`

	// OnCompletion (apan) defines the what additional actions occur with the
	// transaction.
	// Valid types:
	// * noop
	// * optin
	// * closeout
	// * clear
	// * update
	// * update
	// * delete
	OnCompletion string `json:"on-completion,omitempty"`
}

type ApplicationResponse struct {
	// Application index and its parameters
	Application Application `json:"application,omitempty"`

	// CurrentRound round at which the results were computed.
	CurrentRound uint64 `json:"current-round,omitempty"`
}

type ApplicationsResponse struct {
	Applications []Application `json:"applications,omitempty"`

	// CurrentRound round at which the results were computed.
	CurrentRound uint64 `json:"current-round,omitempty"`

	// NextToken used for pagination, when making another request provide this token
	// with the next parameter.
	NextToken string `json:"next-token,omitempty"`
}
