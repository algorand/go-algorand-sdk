package models

// TransactionApplication fields for application transactions.
// Definition:
// data/transactions/application.go : ApplicationCallTxnFields
type TransactionApplication struct {
	// Accounts (apat) List of accounts in addition to the sender that may be accessed
	// from the application's approval-program and clear-state-program.
	Accounts []string `json:"accounts,omitempty"`

	// ApplicationArgs (apaa) transaction specific arguments accessed from the
	// application's approval-program and clear-state-program.
	ApplicationArgs [][]byte `json:"application-args,omitempty"`

	// ApplicationId (apid) ID of the application being configured or empty if
	// creating.
	ApplicationId uint64 `json:"application-id"`

	// ApprovalProgram (apap) Logic executed for every application transaction, except
	// when on-completion is set to "clear". It can read and write global state for the
	// application, as well as account-specific local state. Approval programs may
	// reject the transaction.
	ApprovalProgram []byte `json:"approval-program,omitempty"`

	// ClearStateProgram (apsu) Logic executed for application transactions with
	// on-completion set to "clear". It can read and write global state for the
	// application, as well as account-specific local state. Clear state programs
	// cannot reject the transaction.
	ClearStateProgram []byte `json:"clear-state-program,omitempty"`

	// ExtraProgramPages (epp) specifies the additional app program len requested in
	// pages.
	ExtraProgramPages uint64 `json:"extra-program-pages,omitempty"`

	// ForeignApps (apfa) Lists the applications in addition to the application-id
	// whose global states may be accessed by this application's approval-program and
	// clear-state-program. The access is read-only.
	ForeignApps []uint64 `json:"foreign-apps,omitempty"`

	// ForeignAssets (apas) lists the assets whose parameters may be accessed by this
	// application's ApprovalProgram and ClearStateProgram. The access is read-only.
	ForeignAssets []uint64 `json:"foreign-assets,omitempty"`

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
