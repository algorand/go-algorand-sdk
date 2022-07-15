package models

// BoxesResponse box names of an application
type BoxesResponse struct {
	// Boxes
	Boxes []BoxDescriptor `json:"boxes"`
}
