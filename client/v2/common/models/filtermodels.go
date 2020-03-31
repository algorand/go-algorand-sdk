package models

import "time"

// GetPendingTransactionsByAddressParams defines parameters for GetPendingTransactionsByAddress.
type GetPendingTransactionsByAddressParams struct {
	// Truncated number of transactions to display. If max=0, returns all pending txns.
	Max uint64 `url:"max,omitempty"`
	// Return raw msgpack block bytes or json
	format string `url:"format,omitempty"`
}

func NewPendingTransactionsByAddressParams() GetPendingTransactionsByAddressParams {
	// this SDK only uses format msgpack
	return &GetPendingTransactionsByAddressParams{format: "msgpack"}
}

func (p *GetPendingTransactionsByAddressParams) SetMax(max uint64) {
	p.Max = max
}

// GetBlockParams defines parameters for GetBlock.
type GetBlockParams struct {
	// Return raw msgpack block bytes or json
	format string `url:"format,omitempty"`
}

func NewBlockParams() *GetBlockParams {
	return &GetBlockParams{format: "msgpack"}
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

func NewRegisterParticipationKeysParams() *RegisterParticipationKeysAccountIdParams {
	return &RegisterParticipationKeysAccountIdParams{}
}

func (p *RegisterParticipationKeysAccountIdParams) SetFee(fee uint64) {
	p.Fee = fee
}
func (p *RegisterParticipationKeysAccountIdParams) SetDilution(dil uint64) {
	p.KeyDilution = dil
}
func (p *RegisterParticipationKeysAccountIdParams) SetRoundLastValid(rnd uint64) {
	p.RoundLastValid = rnd
}
func (p *RegisterParticipationKeysAccountIdParams) SetNoWait(nowait bool) {
	p.NoWait = nowait
}

// ShutdownParams defines parameters for GetV2Shutdown.
type ShutdownParams struct {
	Timeout uint64 `url:"timeout,omitempty"`
}

func NewShutdownParams() *ShutdownParams {
	return &ShutdownParams{}
}

func (p *ShutdownParams) SetTimeout(timeout uint64) {
	p.Timeout = timeout
}

// GetPendingTransactionsParams defines parameters for GetPendingTransactions.
type GetPendingTransactionsParams struct {
	// Truncated number of transactions to display. If max=0, returns all pending txns.
	Max uint64 `url:"max,omitempty"`

	// Return raw msgpack block bytes or json
	format string `url:"format,omitempty"`
}

func NewPendingTransactionsParams() *GetPendingTransactionsParams {
	return &GetPendingTransactionsParams{format: "msgpack"}
}

func (p *GetPendingTransactionsParams) SetMax(max uint64) {
	p.Max = max
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

func NewSearchAccountsParams() *SearchAccountsParams {
	return &SearchAccountsParams{}
}

func (p *SearchAccountsParams) SetAssetID(assetID uint64) {
	p.AssetId = assetID
}

func (p *SearchAccountsParams) SetLimit(limit uint64) {
	p.Limit = limit
}

func (p *SearchAccountsParams) SetCurrencyGreaterThan(greaterThan uint64) {
	p.CurrencyGreaterThan = greaterThan
}

func (p *SearchAccountsParams) SetCurrencyLessThan(lessThan uint64) {
	p.CurrencyLessThan = lessThan
}

func (p *SearchAccountsParams) SetAfterAddress(after string) {
	p.AfterAddress = after
}

// LookupAccountByIDParams defines parameters for LookupAccountByID.
type LookupAccountByIDParams struct {

	// Include results for the specified round.
	Round uint64 `url:"round,omitempty"`
}

func NewLookupAccountByIDParams() *LookupAccountByIDParams {
	return &LookupAccountByIDParams{}
}

func (p *LookupAccountByIDParams) SetRound(rnd uint64) {
	p.Round = rnd
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

func NewLookupAccountTransactionsParams() *LookupAccountTransactionsParams {
	return &LookupAccountTransactionsParams{}
}

func (p *LookupAccountTransactionsParams) SetNotePrefix(prefix []byte) {
	p.NotePrefix = prefix
}

func (p *LookupAccountTransactionsParams) SetTxType(txtype string) {
	p.TxType = txtype
}

func (p *LookupAccountTransactionsParams) SetSigType(sigtype string) {
	p.SigType = sigtype
}

func (p *LookupAccountTransactionsParams) SetTXID(txid string) {
	p.TxId = txid
}

func (p *LookupAccountTransactionsParams) SetRound(rnd uint64) {
	p.Round = rnd
}

func (p *LookupAccountTransactionsParams) SetMinRound(min uint64) {
	p.MinRound = min
}

func (p *LookupAccountTransactionsParams) SetMaxRound(max uint64) {
	p.MaxRound = max
}

func (p *LookupAccountTransactionsParams) SetAssetID(id uint64) {
	p.AssetId = id
}

func (p *LookupAccountTransactionsParams) SetLimit(limit uint64) {
	p.Limit = limit
}

func (p *LookupAccountTransactionsParams) SetBeforeTime(before time.Time) {
	p.BeforeTime = before
}

func (p *LookupAccountTransactionsParams) SetAfterTime(after time.Time) {
	p.AfterTime = after
}

func (p *LookupAccountTransactionsParams) SetCurrencyGreaterThan(greaterThan uint64) {
	p.CurrencyGreaterThan = greaterThan
}

func (p *LookupAccountTransactionsParams) SetCurrencyLessThan(lessThan uint64) {
	p.CurrencyLessThan = lessThan
}

func (p *LookupAccountTransactionsParams) SetAddressRole(role string) {
	p.AddressRole = role
}

func (p *LookupAccountTransactionsParams) SetExcludeCloseTo(exclude bool) {
	p.ExcludeCloseTo = exclude
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

func NewSearchForAssetsParams() *SearchForAssetsParams {
	return &SearchForAssetsParams{}
}

func (p *SearchForAssetsParams) SetLimit(lim uint64) {
	p.Limit = lim
}

func (p *SearchForAssetsParams) SetCreator(creator string) {
	p.Creator = creator
}

func (p *SearchForAssetsParams) SetName(name string) {
	p.Name = name
}

func (p *SearchForAssetsParams) SetUnit(unit string) {
	p.Unit = unit
}

func (p *SearchForAssetsParams) SetAssetID(id uint64) {
	p.AssetId = id
}

func (p *SearchForAssetsParams) SetAfterAsset(after uint64) {
	p.AfterAsset = after
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

func NewLookupAssetBalancesParams() *LookupAssetBalancesParams {
	return &LookupAssetBalancesParams{}
}

func (p *LookupAssetBalancesParams) SetLimit(lim uint64) {
	p.Limit = lim
}

func (p *LookupAssetBalancesParams) SetAfterAddress(after string) {
	p.AfterAddress = after
}

func (p *LookupAssetBalancesParams) SetRound(rnd uint64) {
	p.Round = rnd
}

func (p *LookupAssetBalancesParams) SetCurrencyGreaterThan(greaterThan uint64) {
	p.CurrencyGreaterThan = greaterThan
}

func (p *LookupAssetBalancesParams) SetCurrencyLessThan(lessThan uint64) {
	p.CurrencyLessThan = lessThan
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

func NewLookupAssetTransactionsParams() *LookupAssetTransactionsParams {
	return &LookupAssetTransactionsParams{}
}

func (p *LookupAssetTransactionsParams) SetNotePrefix(prefix []byte) {
	p.NotePrefix = prefix
}

func (p *LookupAssetTransactionsParams) SetTxType(txtype string) {
	p.TxType = txtype
}

func (p *LookupAssetTransactionsParams) SetSigType(sigtype string) {
	p.SigType = sigtype
}

func (p *LookupAssetTransactionsParams) SetTXID(txid string) {
	p.TxId = txid
}

func (p *LookupAssetTransactionsParams) SetRound(rnd uint64) {
	p.Round = rnd
}

func (p *LookupAssetTransactionsParams) SetMinRound(min uint64) {
	p.MinRound = min
}

func (p *LookupAssetTransactionsParams) SetMaxRound(max uint64) {
	p.MaxRound = max
}

func (p *LookupAssetTransactionsParams) SetAddress(address string) {
	p.Address = address
}

func (p *LookupAssetTransactionsParams) SetLimit(limit uint64) {
	p.Limit = limit
}

func (p *LookupAssetTransactionsParams) SetBeforeTime(before time.Time) {
	p.BeforeTime = before
}

func (p *LookupAssetTransactionsParams) SetAfterTime(after time.Time) {
	p.AfterTime = after
}

func (p *LookupAssetTransactionsParams) SetCurrencyGreaterThan(greaterThan uint64) {
	p.CurrencyGreaterThan = greaterThan
}

func (p *LookupAssetTransactionsParams) SetCurrencyLessThan(lessThan uint64) {
	p.CurrencyLessThan = lessThan
}

func (p *LookupAssetTransactionsParams) SetAddressRole(role string) {
	p.AddressRole = role
}

func (p *LookupAssetTransactionsParams) SetExcludeCloseTo(exclude bool) {
	p.ExcludeCloseTo = exclude
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

func NewSearchForTransactionsParams() *SearchForTransactionsParams {
	return &SearchForTransactionsParams{}
}

func (p *SearchForTransactionsParams) SetNotePrefix(prefix []byte) {
	p.NotePrefix = prefix
}

func (p *SearchForTransactionsParams) SetTxType(txtype string) {
	p.TxType = txtype
}

func (p *SearchForTransactionsParams) SetSigType(sigtype string) {
	p.SigType = sigtype
}

func (p *SearchForTransactionsParams) SetTXID(txid string) {
	p.TxId = txid
}

func (p *SearchForTransactionsParams) SetRound(rnd uint64) {
	p.Round = rnd
}

func (p *SearchForTransactionsParams) SetMinRound(min uint64) {
	p.MinRound = min
}

func (p *SearchForTransactionsParams) SetMaxRound(max uint64) {
	p.MaxRound = max
}

func (p *SearchForTransactionsParams) SetAddress(address string) {
	p.Address = address
}

func (p *SearchForTransactionsParams) SetLimit(limit uint64) {
	p.Limit = limit
}

func (p *SearchForTransactionsParams) SetBeforeTime(before time.Time) {
	p.BeforeTime = before
}

func (p *SearchForTransactionsParams) SetAfterTime(after time.Time) {
	p.AfterTime = after
}

func (p *SearchForTransactionsParams) SetCurrencyGreaterThan(greaterThan uint64) {
	p.CurrencyGreaterThan = greaterThan
}

func (p *SearchForTransactionsParams) SetCurrencyLessThan(lessThan uint64) {
	p.CurrencyLessThan = lessThan
}

func (p *SearchForTransactionsParams) SetAddressRole(role string) {
	p.AddressRole = role
}

func (p *SearchForTransactionsParams) SetExcludeCloseTo(exclude bool) {
	p.ExcludeCloseTo = exclude
}

func (p *SearchForTransactionsParams) SetAssetID(id uint64) {
	p.AssetId = id
}
