package models

// AccountApplicationResponse defines model for AccountApplicationResponse.
type AccountApplicationResponse struct {
	// Stores local state associated with an application.
	AppLocalState *ApplicationLocalState `json:"app-local-state,omitempty"`

	// Stores the global information associated with an application.
	CreatedApp *ApplicationParams `json:"created-app,omitempty"`

	// The round for which this information is relevant.
	Round uint64 `json:"round"`
}
