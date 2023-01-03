package models

// ModifiedAsset asset which was created or deleted.
type ModifiedAsset struct {
	// Created created if true, deleted if false
	Created bool `json:"created"`

	// Creator address of the creator.
	Creator string `json:"creator"`

	// Id asset Id
	Id uint64 `json:"id"`
}
