package models

// BoxReference references a box of an application.
type BoxReference struct {
	// App application ID which this box belongs to
	App uint64 `json:"app"`

	// Name base64 encoded box name
	Name []byte `json:"name"`
}
