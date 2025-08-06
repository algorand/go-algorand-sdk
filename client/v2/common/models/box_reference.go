package models

// BoxReference boxReference names a box by its name and the application ID it
// belongs to.
type BoxReference struct {
	// App application ID to which the box belongs, or zero if referring to the called
	// application.
	App uint64 `json:"app"`

	// Name base64 encoded box name
	Name []byte `json:"name"`
}
