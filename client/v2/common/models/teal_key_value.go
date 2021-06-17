package models

// TealKeyValue represents a key-value pair in an application store.
type TealKeyValue struct {
	// Key
	Key string `json:"key"`

	// Value represents a TEAL value.
	Value TealValue `json:"value"`
}
