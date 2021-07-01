package models

// BuildVersion defines a model for BuildVersion.
type BuildVersion struct {
	// Branch
	Branch string `json:"branch"`

	// BuildNumber
	BuildNumber uint64 `json:"build_number"`

	// Channel
	Channel string `json:"channel"`

	// CommitHash
	CommitHash string `json:"commit_hash"`

	// Major
	Major uint64 `json:"major"`

	// Minor
	Minor uint64 `json:"minor"`
}
