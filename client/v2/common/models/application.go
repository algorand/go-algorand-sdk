package models

// Application application index and its parameters
type Application struct {
	// Id (appidx) application index.
	Id uint64 `json:"id"`

	// Params (appparams) application parameters.
	Params ApplicationParams `json:"params"`
}
