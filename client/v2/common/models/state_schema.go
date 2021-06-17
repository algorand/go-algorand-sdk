package models

// StateSchema represents a (apls) local-state or (apgs) global-state schema. These
// schemas determine how much storage may be used in a local-state or global-state
// for an application. The more space used, the larger minimum balance must be
// maintained in the account holding the data.
type StateSchema struct {
	// NumByteSlice maximum number of TEAL byte slices that may be stored in the
	// key/value store.
	NumByteSlice uint64 `json:"num-byte-slice"`

	// NumUint maximum number of TEAL uints that may be stored in the key/value store.
	NumUint uint64 `json:"num-uint"`
}
