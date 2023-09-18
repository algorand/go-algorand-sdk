package models

// ApplicationKVStorage an application's global/local/box state.
type ApplicationKVStorage struct {
	// Account the address of the account associated with the local state.
	Account string `json:"account,omitempty"`

	// Kvs key-Value pairs representing application states.
	Kvs []AvmKeyValue `json:"kvs"`
}
