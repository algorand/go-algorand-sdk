package models

// BoxDescriptor box descriptor describes an app box without a value.
type BoxDescriptor struct {
	// Name base64 encoded box name
	Name []byte `json:"name"`
}
