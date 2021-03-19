package models

// BuildVersion defines a model for BuildVersion.
type BuildVersion struct {
	// Branch
	Branch string `json:"branch,omitempty"`

	// BuildNumber
	BuildNumber uint64 `json:"build_number,omitempty"`

	// Channel
	Channel string `json:"channel,omitempty"`

	// CommitHash
	CommitHash string `json:"commit_hash,omitempty"`

	// Major
	Major uint64 `json:"major,omitempty"`

	// Minor
	Minor uint64 `json:"minor,omitempty"`
}
