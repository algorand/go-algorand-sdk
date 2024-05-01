package models

// ApplicationStateSchema specifies maximums on the number of each type that may be
// stored.
type ApplicationStateSchema struct {
	// NumByteSlice number of byte slices.
	NumByteSlice uint64 `json:"num-byte-slice"`

	// NumUint number of uints.
	NumUint uint64 `json:"num-uint"`
}
