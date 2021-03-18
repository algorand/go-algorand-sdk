package models

type RawBlockJson struct {
	Block string
}
type RawBlockMsgpack struct {
	Block string `json:"url,omitempty"`
}

// Supply defines model for Supply.
type Supply struct {

	// OnlineMoney
	OnlineMoney uint64 `json:"online-money"`

	// Round
	Round uint64 `json:"current_round"`

	// TotalMoney
	TotalMoney uint64 `json:"total-money"`
}

// TransactionParams contains the parameters that help a client construct a new transaction.
type TransactionParams struct {

	// ConsensusVersion indicates the consensus protocol version
	// as of LastRound.
	ConsensusVersion string `json:"consensus-version"`

	// Fee is the suggested transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but transactions must still have a fee of
	// at least MinTxnFee for the current network protocol.
	Fee uint64 `json:"fee"`

	// GenesisID is an ID listed in the genesis block.
	GenesisID string `json:"genesis-id"`

	// GenesisHash is the hash of the genesis block.
	Genesishash []byte `json:"genesis-hash"`

	// LastRound indicates the last round seen
	LastRound uint64 `json:"last-round"`

	// The minimum transaction fee (not per byte) required for the
	// txn to validate for the current network protocol.
	MinFee uint64 `json:"min-fee,omitempty"`
}

// VersionBuild defines model for the current algod build version information.
type VersionBuild struct {
	Branch      string `json:"branch"`
	BuildNumber uint64 `json:"build-number"`
	Channel     string `json:"channel"`
	CommitHash  []byte `json:"commit-hash"`
	Major       uint64 `json:"major"`
	Minor       uint64 `json:"minor"`
}

// AccountId defines model for account-id.
type AccountId string

// Address defines model for address.
type Address string

// AddressGreaterThan defines model for address-greater-than.
type AddressGreaterThan string

// AddressRole defines model for address-role.
type AddressRole string

// AfterAddress defines model for after-address.
type AfterAddress string

// AfterAsset defines model for after-asset.
type AfterAsset uint64

// AfterTime defines model for after-time.
type AfterTime string

// AlgosGreaterThan defines model for algos-greater-than.
type AlgosGreaterThan uint64

// AlgosLessThan defines model for algos-less-than.
type AlgosLessThan uint64

// AssetId defines model for asset-id.
type AssetId uint64

// BeforeTime defines model for before-time.
type BeforeTime string

// CurrencyGreaterThan defines model for currency-greater-than.
type CurrencyGreaterThan uint64

// CurrencyLessThan defines model for currency-less-than.
type CurrencyLessThan uint64

// ExcludeCloseTo defines model for exclude-close-to.
type ExcludeCloseTo bool

// Limit defines model for limit.
type Limit uint64

// MaxRound defines model for max-round.
type MaxRound uint64

// MinRound defines model for min-round.
type MinRound uint64

// NotePrefix defines model for note-prefix.
type NotePrefix []byte

// Offset defines model for offset.
type Offset uint64

// Round defines model for round.
type Round uint64

// RoundNumber defines model for round-number.
type RoundNumber uint64

// SigType defines model for sig-type.
type SigType string

// TxId defines model for tx-id.
type TxId []byte

// TxType defines model for tx-type.
type TxType string

// HealthCheckResponse defines model for HealthCheckResponse.
type HealthCheckResponse HealthCheck

// GetBlock response is returned by Block
type GetBlockResponse = struct {
	Blockb64 string `json:"block"`
}

type LookupAccountByIDResponse struct {
	CurrentRound uint64  `json:"current-round"`
	Account      Account `json:"account"`
}

type LookupAssetByIDResponse struct {
	CurrentRound uint64 `json:"current-round"`
	Asset        Asset  `json:"asset"`
}
