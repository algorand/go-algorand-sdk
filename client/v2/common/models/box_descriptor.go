package models

// BoxDescriptor box descriptor describes an app box.
type BoxDescriptor struct {
	// Name base64 encoded box name
	Name []byte `json:"name"`
	// Value base64 encoded box value. Present only when the `values` query parameter is set to true.
	Value []byte `json:"value,omitempty"`
}
