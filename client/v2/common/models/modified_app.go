package models

// ModifiedApp app which was created or deleted.
type ModifiedApp struct {
	// Created created if true, deleted if false
	Created bool `json:"created"`

	// Creator address of the creator.
	Creator string `json:"creator"`

	// Id app Id
	Id uint64 `json:"id"`
}
