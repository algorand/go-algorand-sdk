package models

// ApplicationLocalState stores local state associated with an application.
type ApplicationLocalState struct {
	// Id the application which this local state is for.
	Id uint64 `json:"id"`

	// KeyValue (tkv) storage.
	KeyValue []TealKeyValue `json:"key-value,omitempty"`

	// Schema (hsch) schema.
	Schema ApplicationStateSchema `json:"schema"`
}
