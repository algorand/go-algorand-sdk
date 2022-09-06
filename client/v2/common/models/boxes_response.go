package models

// BoxesResponse box names of an application
type BoxesResponse struct {
	// ApplicationId (appidx) application index.
	ApplicationId uint64 `json:"application-id"`

	// Boxes
	Boxes []BoxDescriptor `json:"boxes"`
}
