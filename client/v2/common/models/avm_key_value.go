package models

// AvmKeyValue represents an AVM key-value pair in an application store.
type AvmKeyValue struct {
	// Key
	Key []byte `json:"key"`

	// Value represents an AVM value.
	Value AvmValue `json:"value"`
}
