package models

// ApplicationLocalState stores local state associated with an application.
type ApplicationLocalState struct {
	// ClosedOutAtRound round when account closed out of the application.
	ClosedOutAtRound uint64 `json:"closed-out-at-round,omitempty"`

	// Deleted whether or not the application local state is currently deleted from its
	// account.
	Deleted bool `json:"deleted,omitempty"`

	// Id the application which this local state is for.
	Id uint64 `json:"id"`

	// KeyValue (tkv) storage.
	KeyValue []TealKeyValue `json:"key-value,omitempty"`

	// OptedInAtRound round when the account opted into the application.
	OptedInAtRound uint64 `json:"opted-in-at-round,omitempty"`

	// Schema (hsch) schema.
	Schema ApplicationStateSchema `json:"schema"`
}
