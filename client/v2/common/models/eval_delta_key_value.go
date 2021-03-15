package models;

// EvalDeltaKeyValue key-value pairs for StateDelta.
type EvalDeltaKeyValue struct {
   // Key
  Key string `json:"key,omitempty"`

   // Value represents a TEAL value delta.
  Value EvalDelta `json:"value,omitempty"`
}
