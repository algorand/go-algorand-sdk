package types

// TxType identifies the type of the transaction
type TxType string

const (
	// PaymentTx is the TxType for payment transactions
	PaymentTx TxType = "pay"
	// KeyRegistrationTx is the TxType for key registration transactions
	KeyRegistrationTx TxType = "keyreg"
	// AssetConfigTx creates, re-configures, or destroys an asset
	AssetConfigTx TxType = "acfg"
	// AssetTransferTx transfers assets between accounts (optionally closing)
	AssetTransferTx TxType = "axfer"
	// AssetFreezeTx changes the freeze status of an asset
	AssetFreezeTx TxType = "afrz"
)

// Algos are the unit of currency in Algorand
type Algos uint64

// Round represents a round of the Algorand consensus protocol
type Round uint64

// VotePK is the participation public key used in key registration transactions
type VotePK [32]byte

// VRFPK is the VRF public key used in key registration transactions
type VRFPK [32]byte

// MasterDerivationKey is the secret key used to derive keys in wallets
type MasterDerivationKey [32]byte

// Digest is a SHA512_256 hash
type Digest [hashLenBytes]byte
