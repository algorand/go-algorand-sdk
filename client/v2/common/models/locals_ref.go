package models

// LocalsRef localsRef names a local state by referring to an Address and App it
// belongs to.
type LocalsRef struct {
	// Address (d) Address in access list, or the sender of the transaction.
	Address string `json:"address"`

	// App (p) Application ID for app in access list, or zero if referring to the
	// called application.
	App uint64 `json:"app"`
}
