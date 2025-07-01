package models

// BoxReference boxReference names a box by its name and the application ID it
// belongs to.
type BoxReference struct {
	// ApplicationId application ID to which the box belongs, or zero if referring to
	// the called application.
	ApplicationId uint64 `json:"application-id"`

	// Name base64 encoded box name
	Name []byte `json:"name"`
}
