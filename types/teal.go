package types

// DeltaAction is an enum of actions that may be performed when applying a
// delta to a TEAL key/value store
type DeltaAction uint64

const (
	// SetUintAction indicates that a Uint should be stored at a key
	SetUintAction DeltaAction = 1

	// SetBytesAction indicates that a TEAL byte slice should be stored at a key
	SetBytesAction DeltaAction = 2

	// DeleteAction indicates that the value for a particular key should be deleted
	DeleteAction DeltaAction = 3
)

// ValueDelta links a DeltaAction with a value to be set
type ValueDelta struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Action DeltaAction `codec:"at"`
	Bytes  string      `codec:"bs"`
	Uint   uint64      `codec:"ui"`
}

// ToTealValue converts a ValueDelta into a TealValue if possible, and returns
// ok = false with a value of 0 TealUint if the conversion is not possible.
func (vd *ValueDelta) ToTealValue() (value TealValue, ok bool) {
	switch vd.Action {
	case SetBytesAction:
		value.Type = TealBytesType
		value.Bytes = vd.Bytes
		ok = true
	case SetUintAction:
		value.Type = TealUintType
		value.Uint = vd.Uint
		ok = true
	case DeleteAction:
		value.Type = TealUintType
		ok = false
	default:
		value.Type = TealUintType
		ok = false
	}
	return value, ok
}

// StateDelta is a map from key/value store keys to ValueDeltas, indicating
// what should happen for that key
//msgp:allocbound StateDelta config.MaxStateDeltaKeys
type StateDelta map[string]ValueDelta

// Equal checks whether two StateDeltas are equal. We don't check for nilness
// equality because an empty map will encode/decode as nil. So if our generated
// map is empty but not nil, we want to equal a decoded nil off the wire.
func (sd StateDelta) Equal(o StateDelta) bool {
	// Lengths should be the same
	if len(sd) != len(o) {
		return false
	}
	// All keys and deltas should be the same
	for k, v := range sd {
		// Other StateDelta must contain key
		ov, ok := o[k]
		if !ok {
			return false
		}

		// Other StateDelta must have same value for key
		if ov != v {
			return false
		}
	}
	return true
}

// EvalDelta stores StateDeltas for an application's global key/value store, as
// well as StateDeltas for some number of accounts holding local state for that
// application
type EvalDelta struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	GlobalDelta StateDelta `codec:"gd"`

	// When decoding EvalDeltas, the integer key represents an offset into
	// [txn.Sender, txn.Accounts[0], txn.Accounts[1], ...]
	LocalDeltas map[uint64]StateDelta `codec:"ld,allocbound=config.MaxEvalDeltaAccounts"`
}

// Equal compares two EvalDeltas and returns whether or not they are
// equivalent. It does not care about nilness equality of LocalDeltas,
// because the msgpack codec will encode/decode an empty map as nil, and we want
// an empty generated EvalDelta to equal an empty one we decode off the wire.
func (ed EvalDelta) Equal(o EvalDelta) bool {
	// LocalDeltas length should be the same
	if len(ed.LocalDeltas) != len(o.LocalDeltas) {
		return false
	}

	// All keys and local StateDeltas should be the same
	for k, v := range ed.LocalDeltas {
		// Other LocalDelta must have value for key
		ov, ok := o.LocalDeltas[k]
		if !ok {
			return false
		}

		// Other LocalDelta must have same value for key
		if !ov.Equal(v) {
			return false
		}
	}

	// GlobalDeltas must be equal
	if !ed.GlobalDelta.Equal(o.GlobalDelta) {
		return false
	}

	return true
}

// StateSchema sets maximums on the number of each type that may be stored
type StateSchema struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	NumUint      uint64 `codec:"nui"`
	NumByteSlice uint64 `codec:"nbs"`
}


// TealType is an enum of the types in a TEAL program: Bytes and Uint
type TealType uint64

const (
	// TealBytesType represents the type of a byte slice in a TEAL program
	TealBytesType TealType = 1

	// TealUintType represents the type of a uint in a TEAL program
	TealUintType TealType = 2
)


// TealValue contains type information and a value, representing a value in a
// TEAL program
type TealValue struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Type  TealType `codec:"tt"`
	Bytes string   `codec:"tb"`
	Uint  uint64   `codec:"ui"`
}

// TealKeyValue represents a key/value store for use in an application's
// LocalState or GlobalState
//msgp:allocbound TealKeyValue 4096
type TealKeyValue map[string]TealValue

