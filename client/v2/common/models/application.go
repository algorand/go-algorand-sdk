package models

// Application application index and its parameters
type Application struct {
	// CreatedAtRound round when this application was created.
	CreatedAtRound uint64 `json:"created-at-round,omitempty"`

	// Deleted whether or not this application is currently deleted.
	Deleted bool `json:"deleted,omitempty"`

	// DeletedAtRound round when this application was deleted.
	DeletedAtRound uint64 `json:"deleted-at-round,omitempty"`

	// Id (appidx) application index.
	Id uint64 `json:"id"`

	// Params (appparams) application parameters.
	Params ApplicationParams `json:"params"`
}
