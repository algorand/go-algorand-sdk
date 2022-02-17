package models

// AccountApplicationResponse accountApplicationResponse describes the application
// data for a specific account and application ID.
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
