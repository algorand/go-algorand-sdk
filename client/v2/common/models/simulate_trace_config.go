package models

// SimulateTraceConfig an object that configures simulation execution trace.
type SimulateTraceConfig struct {
	// Enable a boolean option for opting in execution trace features simulation
	// endpoint.
	Enable bool `json:"enable,omitempty"`
}
