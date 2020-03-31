package models

import "time"

// GetPendingTransactionsByAddressParams defines parameters for GetPendingTransactionsByAddress.
type GetPendingTransactionsByAddressParams struct {
	// Truncated number of transactions to display. If max=0, returns all pending txns.
	Max uint64 `url:"max,omitempty"`
	// Return raw msgpack block bytes or json
	format string `url:"format,omitempty"`
}

type PendingTransactionsByAddressOption func(GetPendingTransactionsByAddressParams) GetPendingTransactionsByAddressParams

func NewPendingTransactionsByAddressParams(options ...PendingTransactionsByAddressOption) GetPendingTransactionsByAddressParams {
	// this SDK only uses msgpack format.
	p := GetPendingTransactionsByAddressParams{format: "msgpack"}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func MaxTransactionsByAddr(max uint64) PendingTransactionsByAddressOption {
	return func(p GetPendingTransactionsByAddressParams) GetPendingTransactionsByAddressParams {
		p.Max = max
		return p
	}
}

// GetBlockParams defines parameters for GetBlock.
type GetBlockParams struct {
	// Return raw msgpack block bytes or json
	format string `url:"format,omitempty"`
}

type GetBlockOption func(GetBlockParams) GetBlockParams

func NewBlockParams(options ...GetBlockOption) GetBlockParams {
	// this SDK only uses msgpack format.
	p := GetBlockParams{format: "msgpack"}
	for _, option := range options {
		p = option(p)
	}
	return p
}

// RegisterParticipationKeysAccountIdParams defines parameters for GetV2RegisterParticipationKeysAccountId.
type RegisterParticipationKeysAccountIdParams struct {

	// The fee to use when submitting key registration transactions. Defaults to the suggested fee.
	Fee uint64 `url:"fee,omitempty"`

	// value to use for two-level participation key.
	KeyDilution uint64 `url:"key-dilution,omitempty"`

	// The last round for which the generated participation keys will be valid.
	RoundLastValid uint64 `url:"round-last-valid,omitempty"`

	// Don't wait for transaction to commit.
	NoWait bool `url:"no-wait,omitempty"`
}

type RegisterParticipationKeyOption func(RegisterParticipationKeysAccountIdParams) RegisterParticipationKeysAccountIdParams

func NewRegisterParticipationKeysParams(options ...RegisterParticipationKeyOption) RegisterParticipationKeysAccountIdParams {
	p := RegisterParticipationKeysAccountIdParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func RegisterParticipationFee(fee uint64) RegisterParticipationKeyOption {
	return func(p RegisterParticipationKeysAccountIdParams) RegisterParticipationKeysAccountIdParams {
		p.Fee = fee
		return p
	}
}

func RegisterParticipationDilution(dil uint64) RegisterParticipationKeyOption {
	return func(p RegisterParticipationKeysAccountIdParams) RegisterParticipationKeysAccountIdParams {
		p.KeyDilution = dil
		return p
	}
}

func RegisterParticipationRoundLastValid(rnd uint64) RegisterParticipationKeyOption {
	return func(p RegisterParticipationKeysAccountIdParams) RegisterParticipationKeysAccountIdParams {
		p.RoundLastValid = rnd
		return p
	}
}

func RegisterParticipationNoWait(nowait bool) RegisterParticipationKeyOption {
	return func(p RegisterParticipationKeysAccountIdParams) RegisterParticipationKeysAccountIdParams {
		p.NoWait = nowait
		return p
	}
}

// ShutdownParams defines parameters for GetV2Shutdown.
type ShutdownParams struct {
	Timeout uint64 `url:"timeout,omitempty"`
}

type ShutdownOption func(ShutdownParams) ShutdownParams

func NewShutdownParams(options ...ShutdownOption) ShutdownParams {
	p := ShutdownParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func ShutdownTimeout(timeout uint64) ShutdownOption {
	return func(p ShutdownParams) ShutdownParams {
		p.Timeout = timeout
		return p
	}
}

// GetPendingTransactionsParams defines parameters for GetPendingTransactions.
type GetPendingTransactionsParams struct {
	// Truncated number of transactions to display. If max=0, returns all pending txns.
	Max uint64 `url:"max,omitempty"`

	// Return raw msgpack block bytes or json
	format string `url:"format,omitempty"`
}

type PendingTransactionsOption func(params GetPendingTransactionsParams) GetPendingTransactionsParams

func NewPendingTransactionsParams(options ...PendingTransactionsOption) GetPendingTransactionsParams {
	// this SDK only uses msgpack format.
	p := GetPendingTransactionsParams{format: "msgpack"}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func MaxPendingTransactions(max uint64) PendingTransactionsOption {
	return func(p GetPendingTransactionsParams) GetPendingTransactionsParams {
		p.Max = max
		return p
	}
}

// SearchAccountsParams defines parameters for SearchAccounts.
type SearchAccountsParams struct {

	// Include accounts holding the specified asset
	AssetId uint64 `url:"asset-id,omitempty"`

	// Maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Results should have an amount greater than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	// Results should have an amount less than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	// Used in conjunction with limit to page through results.
	AfterAddress string `url:"after-address,omitempty"`
}

type SearchAccountsOption func(params SearchAccountsParams) SearchAccountsParams

func NewSearchAccountsParams(options ...SearchAccountsOption) SearchAccountsParams {
	p := SearchAccountsParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func SearchAccountsAssetID(assetID uint64) SearchAccountsOption {
	return func(p SearchAccountsParams) SearchAccountsParams {
		p.AssetId = assetID
		return p
	}
}

func SearchAccountsLimit(limit uint64) SearchAccountsOption {
	return func(p SearchAccountsParams) SearchAccountsParams {
		p.Limit = limit
		return p
	}
}

func SearchAccountsCurrencyGreaterThan(greaterThan uint64) SearchAccountsOption {
	return func(p SearchAccountsParams) SearchAccountsParams {
		p.CurrencyGreaterThan = greaterThan
		return p
	}
}

func SearchAccountsCurrencyLessThan(lessThan uint64) SearchAccountsOption {
	return func(p SearchAccountsParams) SearchAccountsParams {
		p.CurrencyLessThan = lessThan
		return p
	}
}

func SearchAccountsAfterAddress(after string) SearchAccountsOption {
	return func(p SearchAccountsParams) SearchAccountsParams {
		p.AfterAddress = after
		return p
	}
}

// LookupAccountByIDParams defines parameters for LookupAccountByID.
type LookupAccountByIDParams struct {

	// Include results for the specified round.
	Round uint64 `url:"round,omitempty"`
}

type LookupAccountByIDOption func(LookupAccountByIDParams) LookupAccountByIDParams

func NewLookupAccountByIDParams(options ...LookupAccountByIDOption) LookupAccountByIDParams {
	p := LookupAccountByIDParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func LookupAccountByIDRound(rnd uint64) LookupAccountByIDOption {
	return func(p LookupAccountByIDParams) LookupAccountByIDParams {
		p.Round = rnd
		return p
	}
}

// LookupAccountTransactionsParams defines parameters for LookupAccountTransactions.
type LookupAccountTransactionsParams struct {

	// Specifies a prefix which must be contained in the note field.
	NotePrefix []byte `url:"note-prefix,omitempty"`
	TxType     string `url:"tx-type,omitempty"`

	// SigType filters just results using the specified type of signature:
	//  sig - Standard
	//  msig - MultiSig
	//  lsig - LogicSig
	SigType string `url:"sig-type,omitempty"`

	// Lookup the specific transaction by ID.
	TxId string `url:"tx-id,omitempty"`

	// Include results for the specified round.
	Round uint64 `url:"round,omitempty"`

	// Include results at or after the specified min-round.
	MinRound uint64 `url:"min-round,omitempty"`

	// Include results at or before the specified max-round.
	MaxRound uint64 `url:"max-round,omitempty"`

	// Asset ID
	AssetId uint64 `url:"asset-id,omitempty"`

	// Maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Include results before the given time. Must be an RFC 3339 formatted string.
	BeforeTime time.Time `url:"before-time,omitempty"`

	// Include results after the given time. Must be an RFC 3339 formatted string.
	AfterTime time.Time `url:"after-time,omitempty"`

	// Results should have an amount greater than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	// Results should have an amount less than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	// Combine with the address parameter to define what type of address to search for.
	AddressRole string `url:"address-role,omitempty"`

	// Combine with address and address-role parameters to define what type of address to search for. The close to fields are normally treated as a receiver, if you would like to exclude them set this parameter to true.
	ExcludeCloseTo bool `url:"exclude-close-to,omitempty"`
}

type LookupAccountTransactionsOption func(LookupAccountTransactionsParams) LookupAccountTransactionsParams

func NewLookupAccountTransactionsParams(options ...LookupAccountTransactionsOption) LookupAccountTransactionsParams {
	p := LookupAccountTransactionsParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func LookupAccountTransactionNotePrefix(prefix []byte) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.NotePrefix = prefix
		return p
	}
}

func LookupAccountTransactionTxType(txtype string) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.TxType = txtype
		return p
	}
}

func LookupAccountTransactionSigType(sigtype string) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.SigType = sigtype
		return p
	}
}

func LookupAccountTransactionTXID(txid string) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.TxId = txid
		return p
	}
}

func LookupAccountTransactionRound(rnd uint64) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.Round = rnd
		return p
	}
}

func LookupAccountTransactionMinRound(min uint64) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.MinRound = min
		return p
	}
}

func LookupAccountTransactionMaxRound(max uint64) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.MaxRound = max
		return p
	}
}

func LookupAccountTransactionAssetID(id uint64) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.AssetId = id
		return p
	}
}

func LookupAccountTransactionLimit(limit uint64) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.Limit = limit
		return p
	}
}

func LookupAccountTransactionBeforeTime(before time.Time) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.BeforeTime = before
		return p
	}
}

func LookupAccountTransactionAfterTime(after time.Time) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.AfterTime = after
		return p
	}
}

func LookupAccountTransactionCurrencyGreaterThan(greaterThan uint64) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.CurrencyGreaterThan = greaterThan
		return p
	}
}

func LookupAccountTransactionCurrencyLessThan(lessThan uint64) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.CurrencyLessThan = lessThan
		return p
	}
}

func LookupAccountTransactionAddressRole(role string) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.AddressRole = role
		return p
	}
}

func LookupAccountTransactionExcludeCloseTo(exclude bool) LookupAccountTransactionsOption {
	return func(p LookupAccountTransactionsParams) LookupAccountTransactionsParams {
		p.ExcludeCloseTo = exclude
		return p
	}
}

// SearchForAssetsParams defines parameters for SearchForAssets.
type SearchForAssetsParams struct {

	// Maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Filter just assets with the given creator address.
	Creator string `url:"creator,omitempty"`

	// Filter just assets with the given name.
	Name string `url:"name,omitempty"`

	// Filter just assets with the given unit.
	Unit string `url:"unit,omitempty"`

	// Asset ID
	AssetId uint64 `url:"asset-id,omitempty"`

	// Used in conjunction with limit to page through results.
	AfterAsset uint64 `url:"after-asset,omitempty"`
}

type SearchForAssetsOption func(SearchForAssetsParams) SearchForAssetsParams

func NewSearchForAssetsParams(options ...SearchForAssetsOption) SearchForAssetsParams {
	p := SearchForAssetsParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func SearchForAssetsLimit(lim uint64) SearchForAssetsOption {
	return func(p SearchForAssetsParams) SearchForAssetsParams {
		p.Limit = lim
		return p
	}
}

func SearchForAssetsCreator(creator string) SearchForAssetsOption {
	return func(p SearchForAssetsParams) SearchForAssetsParams {
		p.Creator = creator
		return p
	}
}

func SearchForAssetsName(name string) SearchForAssetsOption {
	return func(p SearchForAssetsParams) SearchForAssetsParams {
		p.Name = name
		return p
	}
}

func SearchForAssetsUnit(unit string) SearchForAssetsOption {
	return func(p SearchForAssetsParams) SearchForAssetsParams {
		p.Unit = unit
		return p
	}
}

func SearchForAssetsAssetID(id uint64) SearchForAssetsOption {
	return func(p SearchForAssetsParams) SearchForAssetsParams {
		p.AssetId = id
		return p
	}
}

func SearchForAssetsAfterAsset(after uint64) SearchForAssetsOption {
	return func(p SearchForAssetsParams) SearchForAssetsParams {
		p.AfterAsset = after
		return p
	}
}

// LookupAssetBalancesParams defines parameters for LookupAssetBalances.
type LookupAssetBalancesParams struct {

	// Maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Used in conjunction with limit to page through results.
	AfterAddress string `url:"after-address,omitempty"`

	// Include results for the specified round.
	Round uint64 `url:"round,omitempty"`

	// Results should have an amount greater than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	// Results should have an amount less than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`
}

type LookupAssetBalancesOption func(LookupAssetBalancesParams) LookupAssetBalancesParams

func NewLookupAssetBalancesParams(options ...LookupAssetBalancesOption) LookupAssetBalancesParams {
	p := LookupAssetBalancesParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func LookupAssetBalancesLimit(lim uint64) LookupAssetBalancesOption {
	return func(p LookupAssetBalancesParams) LookupAssetBalancesParams {
		p.Limit = lim
		return p
	}
}

func LookupAssetAfterAddress(after string) LookupAssetBalancesOption {
	return func(p LookupAssetBalancesParams) LookupAssetBalancesParams {
		p.AfterAddress = after
		return p
	}
}

func LookupAssetBalancesRound(rnd uint64) LookupAssetBalancesOption {
	return func(p LookupAssetBalancesParams) LookupAssetBalancesParams {
		p.Round = rnd
		return p
	}
}

func LookupAssetBalancesCurrencyGreaterThan(greaterThan uint64) LookupAssetBalancesOption {
	return func(p LookupAssetBalancesParams) LookupAssetBalancesParams {
		p.CurrencyGreaterThan = greaterThan
		return p
	}
}

func LookupAssetBalancesCurrencyLessThan(lessThan uint64) LookupAssetBalancesOption {
	return func(p LookupAssetBalancesParams) LookupAssetBalancesParams {
		p.CurrencyLessThan = lessThan
		return p
	}
}

// LookupAssetTransactionsParams defines parameters for LookupAssetTransactions.
type LookupAssetTransactionsParams struct {

	// Specifies a prefix which must be contained in the note field.
	NotePrefix []byte `url:"note-prefix,omitempty"`
	TxType     string `url:"tx-type,omitempty"`

	// SigType filters just results using the specified type of signature:
	//  sig - Standard
	//  msig - MultiSig
	//  lsig - LogicSig
	SigType string `url:"sig-type,omitempty"`

	// Lookup the specific transaction by ID.
	TxId string `url:"tx-id,omitempty"`

	// Include results for the specified round.
	Round uint64 `url:"round,omitempty"`

	// Include results at or after the specified min-round.
	MinRound uint64 `url:"min-round,omitempty"`

	// Include results at or before the specified max-round.
	MaxRound uint64 `url:"max-round,omitempty"`

	// Maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Include results before the given time. Must be an RFC 3339 formatted string.
	BeforeTime time.Time `url:"before-time,omitempty"`

	// Include results after the given time. Must be an RFC 3339 formatted string.
	AfterTime time.Time `url:"after-time,omitempty"`

	// Results should have an amount greater than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	// Results should have an amount less than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	// Only include transactions with this address in one of the transaction fields.
	Address string `url:"address,omitempty"`

	// Combine with the address parameter to define what type of address to search for.
	AddressRole string `url:"address-role,omitempty"`

	// Combine with address and address-role parameters to define what type of address to search for. The close to fields are normally treated as a receiver, if you would like to exclude them set this parameter to true.
	ExcludeCloseTo bool `url:"exclude-close-to,omitempty"`
}

type LookupAssetTransactionsOption func(LookupAssetTransactionsParams) LookupAssetTransactionsParams

func NewLookupAssetTransactionsParams(options ...LookupAssetTransactionsOption) LookupAssetTransactionsParams {
	p := LookupAssetTransactionsParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func LookupAssetTransactionsNotePrefix(prefix []byte) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.NotePrefix = prefix
		return p
	}
}

func LookupAssetTransactionsTxType(txtype string) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.TxType = txtype
		return p
	}
}

func LookupAssetTransactionsigType(sigtype string) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.SigType = sigtype
		return p
	}
}

func LookupAssetTransactionsTXID(txid string) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.TxId = txid
		return p
	}
}

func LookupAssetTransactionsRound(rnd uint64) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.Round = rnd
		return p
	}
}

func LookupAssetTransactionsMinRound(min uint64) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.MinRound = min
		return p
	}
}

func LookupAssetTransactionsMaxRound(max uint64) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.MaxRound = max
		return p
	}
}

func LookupAssetTransactionsAddress(address string) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.Address = address
		return p
	}
}

func LookupAssetTransactionsLimit(limit uint64) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.Limit = limit
		return p
	}
}

func LookupAssetTransactionsBeforeTime(before time.Time) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.BeforeTime = before
		return p
	}
}

func LookupAssetTransactionsAfterTime(after time.Time) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.AfterTime = after
		return p
	}
}

func LookupAssetTransactionsCurrencyGreaterThan(greaterThan uint64) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.CurrencyGreaterThan = greaterThan
		return p
	}
}

func LookupAssetTransactionsCurrencyLessThan(lessThan uint64) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.CurrencyLessThan = lessThan
		return p
	}
}

func LookupAssetTransactionsAddressRole(role string) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.AddressRole = role
		return p
	}
}

func LookupAssetTransactionsExcludeCloseTo(exclude bool) LookupAssetTransactionsOption {
	return func(p LookupAssetTransactionsParams) LookupAssetTransactionsParams {
		p.ExcludeCloseTo = exclude
		return p
	}
}

// SearchForTransactionsParams defines parameters for SearchForTransactions.
type SearchForTransactionsParams struct {

	// Specifies a prefix which must be contained in the note field.
	NotePrefix []byte `url:"note-prefix,omitempty"`
	TxType     string `url:"tx-type,omitempty"`

	// SigType filters just results using the specified type of signature:
	//  sig - Standard
	//  msig - MultiSig
	//  lsig - LogicSig
	SigType string `url:"sig-type,omitempty"`

	// Lookup the specific transaction by ID.
	TxId string `url:"tx-id,omitempty"`

	// Include results for the specified round.
	Round uint64 `url:"round,omitempty"`

	// Include results at or after the specified min-round.
	MinRound uint64 `url:"min-round,omitempty"`

	// Include results at or before the specified max-round.
	MaxRound uint64 `url:"max-round,omitempty"`

	// Asset ID
	AssetId uint64 `url:"asset-id,omitempty"`

	// Maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Include results before the given time. Must be an RFC 3339 formatted string.
	BeforeTime time.Time `url:"before-time,omitempty"`

	// Include results after the given time. Must be an RFC 3339 formatted string.
	AfterTime time.Time `url:"after-time,omitempty"`

	// Results should have an amount greater than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	// Results should have an amount less than this value. MicroAlgos are the default currency unless an asset-id is provided, in which case the asset will be used.
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	// Only include transactions with this address in one of the transaction fields.
	Address string `url:"address,omitempty"`

	// Combine with the address parameter to define what type of address to search for.
	AddressRole string `url:"address-role,omitempty"`

	// Combine with address and address-role parameters to define what type of address to search for. The close to fields are normally treated as a receiver, if you would like to exclude them set this parameter to true.
	ExcludeCloseTo bool `url:"exclude-close-to,omitempty"`
}

type SearchForTransactionsOption func(SearchForTransactionsParams) SearchForTransactionsParams

func NewSearchForTransactionsParams(options ...SearchForTransactionsOption) SearchForTransactionsParams {
	p := SearchForTransactionsParams{}
	for _, option := range options {
		p = option(p)
	}
	return p
}

func SearchForTransactionsNotePrefix(prefix []byte) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.NotePrefix = prefix
		return p
	}
}

func SearchForTransactionsTxType(txtype string) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.TxType = txtype
		return p
	}
}

func SearchForTransactionsigType(sigtype string) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.SigType = sigtype
		return p
	}
}

func SearchForTransactionsTXID(txid string) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.TxId = txid
		return p
	}
}

func SearchForTransactionsRound(rnd uint64) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.Round = rnd
		return p
	}
}

func SearchForTransactionsMinRound(min uint64) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.MinRound = min
		return p
	}
}

func SearchForTransactionsMaxRound(max uint64) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.MaxRound = max
		return p
	}
}

func SearchForTransactionsAddress(address string) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.Address = address
		return p
	}
}

func SearchForTransactionsAssetID(assetID uint64) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.AssetId = assetID
		return p
	}
}

func SearchForTransactionsLimit(limit uint64) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.Limit = limit
		return p
	}
}

func SearchForTransactionsBeforeTime(before time.Time) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.BeforeTime = before
		return p
	}
}

func SearchForTransactionsAfterTime(after time.Time) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.AfterTime = after
		return p
	}
}

func SearchForTransactionsCurrencyGreaterThan(greaterThan uint64) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.CurrencyGreaterThan = greaterThan
		return p
	}
}

func SearchForTransactionsCurrencyLessThan(lessThan uint64) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.CurrencyLessThan = lessThan
		return p
	}
}

func SearchForTransactionsAddressRole(role string) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.AddressRole = role
		return p
	}
}

func SearchForTransactionsExcludeCloseTo(exclude bool) SearchForTransactionsOption {
	return func(p SearchForTransactionsParams) SearchForTransactionsParams {
		p.ExcludeCloseTo = exclude
		return p
	}
}
