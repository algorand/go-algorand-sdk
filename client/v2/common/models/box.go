package models

// Box box name and its content.
type Box struct {
	// Name (name) box name, base64 encoded
	Name []byte `json:"name"`

	// Value (value) box value, base64 encoded.
	Value []byte `json:"value"`
}
