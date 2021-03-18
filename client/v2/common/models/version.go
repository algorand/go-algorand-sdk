package models

// Version algod version information.
type Version struct {
	// Build
	Build BuildVersion `json:"build,omitempty"`

	// Genesis_hash_b64
	Genesis_hash_b64 []byte `json:"genesis_hash_b64,omitempty"`

	// Genesis_id
	Genesis_id string `json:"genesis_id,omitempty"`

	// Versions
	Versions []string `json:"versions,omitempty"`
}
