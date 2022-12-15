package models

// AppResourceRecord represents AppParams and AppLocalStateDelta in deltas
type AppResourceRecord struct {
	// Address app account address
	Address string `json:"address"`

	// AppDeleted whether the app was deleted
	AppDeleted bool `json:"app-deleted"`

	// AppIndex app index
	AppIndex uint64 `json:"app-index"`

	// AppLocalState app local state
	AppLocalState ApplicationLocalState `json:"app-local-state,omitempty"`

	// AppLocalStateDeleted whether the app local state was deleted
	AppLocalStateDeleted bool `json:"app-local-state-deleted"`

	// AppParams app params
	AppParams ApplicationParams `json:"app-params,omitempty"`
}
