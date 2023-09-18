package models

// ApplicationLocalReference references an account's local state for an
// application.
type ApplicationLocalReference struct {
	// Account address of the account with the local state.
	Account string `json:"account"`

	// App application ID of the local state application.
	App uint64 `json:"app"`
}
