package models

// BoxDescriptor box descriptor describes a Box.
type BoxDescriptor struct {
	// Name base64 encoded box name
	Name []byte `json:"name"`
}
