package models

// AccountApplicationResource accountApplicationResource describes the account's
// application resource (local state and params if the account is the creator) for
// a specific application ID.
type AccountApplicationResource struct {
	// AppLocalState (appl) the application local data stored in this account.
	// The raw account uses `AppLocalState` for this type.
	AppLocalState ApplicationLocalState `json:"app-local-state,omitempty"`

	// CreatedAtRound round when the account opted into or created the application.
	CreatedAtRound uint64 `json:"created-at-round,omitempty"`

	// Deleted whether the application has been deleted.
	Deleted bool `json:"deleted,omitempty"`

	// Id the application ID.
	Id uint64 `json:"id"`

	// Params (appp) parameters of the application created by this account including
	// app global data.
	// The raw account uses `AppParams` for this type.
	// Only present if the account is the creator and `include=params` is specified.
	Params ApplicationParams `json:"params,omitempty"`
}
