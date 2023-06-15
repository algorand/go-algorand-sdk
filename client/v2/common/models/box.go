package models

// Box box name and its content.
type Box struct {
	// Name (name) box name, base64 encoded
	Name []byte `json:"name"`

	// Round the round for which this information is relevant
	Round uint64 `json:"round"`

	// Value (value) box value, base64 encoded.
	Value []byte `json:"value"`
}
