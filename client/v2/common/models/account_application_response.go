package models

// AccountApplicationResponse accountApplicationResponse describes the account's
// application local state and global state (AppLocalState and AppParams, if either
// exists) for a specific application ID. Global state will only be returned if the
// provided address is the application's creator.
type AccountApplicationResponse struct {
	// AppLocalState (appl) the application local data stored in this account.
	// The raw account uses `AppLocalState` for this type.
	AppLocalState ApplicationLocalState `json:"app-local-state,omitempty"`

	// CreatedApp (appp) parameters of the application created by this account
	// including app global data.
	// The raw account uses `AppParams` for this type.
	CreatedApp ApplicationParams `json:"created-app,omitempty"`

	// Round the round for which this information is relevant.
	Round uint64 `json:"round"`
}
