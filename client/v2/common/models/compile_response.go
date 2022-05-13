package models

// CompileResponse teal compile Result
type CompileResponse struct {
	// Hash base32 SHA512_256 of program bytes (Address style)
	Hash string `json:"hash"`

	// Result base64 encoded program bytes
	Result string `json:"result"`

	// Sourcemap jSON of the source map
	Sourcemap *map[string]interface{} `json:"sourcemap,omitempty"`
}
