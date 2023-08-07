package models

// SimulateTraceConfig an object that configures simulation execution trace.
type SimulateTraceConfig struct {
	// Enable a boolean option for opting in execution trace features simulation
	// endpoint.
	Enable bool `json:"enable,omitempty"`

	// ScratchChange a boolean option enabling returning scratch slot changes together
	// with execution trace during simulation.
	ScratchChange bool `json:"scratch-change,omitempty"`

	// StackChange a boolean option enabling returning stack changes together with
	// execution trace during simulation.
	StackChange bool `json:"stack-change,omitempty"`
}
