package types

// TxType identifies the type of the transaction
type TxType string

const (
	// PaymentTx is the TxType for payment transactions
	PaymentTx TxType = "pay"
	// KeyRegistrationTx is the TxType for key registration transactions
	KeyRegistrationTx TxType = "keyreg"
)

// Algos are the unit of currency in Algorand. We use a struct here to
// discourage using arithmetic directly on the Raw field
type Algos struct {
	Raw uint64
}

// Round represents a round of the Algorand consensus protocol
type Round uint64

// VotePK is the participation public key used in key registration transactions
type VotePK [32]byte

// VRFPK is the VRF public key used in key registration transactions
type VRFPK [32]byte
