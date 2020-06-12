package models

// GetPendingTransactionsByAddressParams defines parameters for GetPendingTransactionsByAddress.
type GetPendingTransactionsByAddressParams struct {
	// Truncated number of transactions to display. If max=0, returns all pending txns.
	Max uint64 `url:"max,omitempty"`
	// Return raw msgpack block bytes or json
	Format string `url:"format,omitempty"`
}

// GetBlockParams defines parameters for GetBlock.
type GetBlockParams struct {
	// Return raw msgpack block bytes or json
	Format string `url:"format,omitempty"`
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

// ShutdownParams defines parameters for GetV2Shutdown.
type ShutdownParams struct {
	Timeout uint64 `url:"timeout,omitempty"`
}

// PendingTransactionInformationParams defines parameters for GetPendingTransactions.
type PendingTransactionInformationParams struct {
	// Return raw msgpack block bytes or json
	Format string `url:"format,omitempty"`

	// Truncated number of transactions to display. If max=0, returns all pending txns.
	Max uint64 `url:"max,omitempty"`
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

	// Used for pagination.
	NextToken string `url:"next,omitempty"`

	// Round for results.
	Round uint64 `url:"round,omitempty"`

	// Include accounts associated with this spending key.
	AuthAddr string `url:"auth-addr,omitempty"`
}

// LookupAccountByIDParams defines parameters for LookupAccountByID.
type LookupAccountByIDParams struct {

	// Include results for the specified round.
	Round uint64 `url:"round,omitempty"`
}


type LookupApplicationParams struct {
	/**
	 * Include results for the specified round.
	 */
	Round uint64 `url:"round,omitempty"`
}

type LookupAccountTransactionsParams struct {
	/**
	 * Include results after the given time. Must be an RFC 3339 formatted string.
	 */
	AfterTime string `url:"after-time,omitempty"`

	/**
	 * Application ID
	 */
	ApplicationId uint64 `url:"application-id,omitempty"`

	/**
	 * Asset ID
	 */
	AssetId uint64 `url:"asset-id,omitempty"`

	/**
	 * Include results before the given time. Must be an RFC 3339 formatted string.
	 */
	BeforeTime string `url:"before-time,omitempty"`

	/**
	 * Results should have an amount greater than this value. MicroAlgos are the
	 * default currency unless an asset-id is provided, in which case the asset will be
	 * used.
	 */
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	/**
	 * Results should have an amount less than this value. MicroAlgos are the default
	 * currency unless an asset-id is provided, in which case the asset will be used.
	 */
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	/**
	 * Maximum number of results to return.
	 */
	Limit uint64 `url:"limit,omitempty"`

	/**
	 * Include results at or before the specified max-round.
	 */
	MaxRound uint64 `url:"max-round,omitempty"`

	/**
	 * Include results at or after the specified min-round.
	 */
	MinRound uint64 `url:"min-round,omitempty"`

	/**
	 * The next page of results. Use the next token provided by the previous results.
	 */
	Next string `url:"next,omitempty"`

	/**
	 * Specifies a prefix which must be contained in the note field.
	 */
	NotePrefix string `url:"note-prefix,omitempty"`

	/**
	 * Include results which include the rekey-to field.
	 */
	RekeyTo bool `url:"rekey-to,omitempty"`

	/**
	 * Include results for the specified round.
	 */
	Round uint64 `url:"round,omitempty"`

	/**
	 * SigType filters just results using the specified type of signature:
	 *   sig - Standard
	 *   msig - MultiSig
	 *   lsig - LogicSig
	 */
	SigType string `url:"sig-type,omitempty"`

	TxType string `url:"tx-type,omitempty"`

	/**
	 * Lookup the specific transaction by ID.
	 */
	Txid string `url:"txid,omitempty"`
}

type SearchForApplicationsParams struct {
	/**
	 * Application ID
	 */
	ApplicationId uint64 `url:"application-id,omitempty"`

	/**
	 * Include results for the specified round.
	 */
	Round uint64 `url:"round,omitempty"`
}

type LookupAssetBalancesParams struct {
	/**
	 * Results should have an amount greater than this value. MicroAlgos are the
	 * default currency unless an asset-id is provided, in which case the asset will be
	 * used.
	 */
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	/**
	 * Results should have an amount less than this value. MicroAlgos are the default
	 * currency unless an asset-id is provided, in which case the asset will be used.
	 */
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	/**
	 * Maximum number of results to return.
	 */
	Limit uint64 `url:"limit,omitempty"`

	/**
	 * The next page of results. Use the next token provided by the previous results.
	 */
	Next string `url:"next,omitempty"`

	/**
	 * Include results for the specified round.
	 */
	Round uint64 `url:"round,omitempty"`
}

type LookupAssetTransactionsParams struct {
	/**
	 * Only include transactions with this address in one of the transaction fields.
	 */
	Address string `url:"address,omitempty"`

	/**
	 * Combine with the address parameter to define what type of address to search for.
	 */
	AddressRole string `url:"address-role,omitempty"`

	/**
	 * Include results after the given time. Must be an RFC 3339 formatted string.
	 */
	AfterTime string `url:"after-time,omitempty"`

	/**
	 * Application ID
	 */
	ApplicationId uint64 `url:"application-id,omitempty"`

	/**
	 * Include results before the given time. Must be an RFC 3339 formatted string.
	 */
	BeforeTime string `url:"before-time,omitempty"`

	/**
	 * Results should have an amount greater than this value. MicroAlgos are the
	 * default currency unless an asset-id is provided, in which case the asset will be
	 * used.
	 */
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	/**
	 * Results should have an amount less than this value. MicroAlgos are the default
	 * currency unless an asset-id is provided, in which case the asset will be used.
	 */
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	/**
	 * Combine with address and address-role parameters to define what type of address
	 * to search for. The close to fields are normally treated as a receiver, if you
	 * would like to exclude them set this parameter to true.
	 */
	ExcludeCloseTo bool `url:"exclude-close-to,omitempty"`

	/**
	 * Maximum number of results to return.
	 */
	Limit uint64 `url:"limit,omitempty"`

	/**
	 * Include results at or before the specified max-round.
	 */
	MaxRound uint64 `url:"max-round,omitempty"`

	/**
	 * Include results at or after the specified min-round.
	 */
	MinRound uint64 `url:"min-round,omitempty"`

	/**
	 * The next page of results. Use the next token provided by the previous results.
	 */
	Next string `url:"next,omitempty"`

	/**
	 * Specifies a prefix which must be contained in the note field.
	 */
	NotePrefix string `url:"note-prefix,omitempty"`

	/**
	 * Include results which include the rekey-to field.
	 */
	RekeyTo bool `url:"rekey-to,omitempty"`

	/**
	 * Include results for the specified round.
	 */
	Round uint64 `url:"round,omitempty"`

	/**
	 * SigType filters just results using the specified type of signature:
	 *   sig - Standard
	 *   msig - MultiSig
	 *   lsig - LogicSig
	 */
	SigType string `url:"sig-type,omitempty"`

	TxType string `url:"tx-type,omitempty"`

	/**
	 * Lookup the specific transaction by ID.
	 */
	Txid string `url:"txid,omitempty"`
}

type SearchForAccountsParams struct {
	/**
	 * Application ID
	 */
	ApplicationId uint64 `url:"application-id,omitempty"`

	/**
	 * Asset ID
	 */
	AssetId uint64 `url:"asset-id,omitempty"`

	/**
	 * Include accounts configured to use this spending key.
	 */
	AuthAddr string `url:"auth-addr,omitempty"`

	/**
	 * Results should have an amount greater than this value. MicroAlgos are the
	 * default currency unless an asset-id is provided, in which case the asset will be
	 * used.
	 */
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	/**
	 * Results should have an amount less than this value. MicroAlgos are the default
	 * currency unless an asset-id is provided, in which case the asset will be used.
	 */
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	/**
	 * Maximum number of results to return.
	 */
	Limit uint64 `url:"limit,omitempty"`

	/**
	 * The next page of results. Use the next token provided by the previous results.
	 */
	Next string `url:"next,omitempty"`

	/**
	 * Include results for the specified round. For performance reasons, this parameter
	 * may be disabled on some configurations.
	 */
	Round uint64 `url:"round,omitempty"`
}

type SearchForAssetsParams struct {
	/**
	 * Asset ID
	 */
	AssetId uint64 `url:"asset-id,omitempty"`

	/**
	 * Filter just assets with the given creator address.
	 */
	Creator string `url:"creator,omitempty"`

	/**
	 * Maximum number of results to return.
	 */
	Limit uint64 `url:"limit,omitempty"`

	/**
	 * Filter just assets with the given name.
	 */
	Name string `url:"name,omitempty"`

	/**
	 * The next page of results. Use the next token provided by the previous results.
	 */
	Next string `url:"next,omitempty"`

	/**
	 * Filter just assets with the given unit.
	 */
	Unit string `url:"unit,omitempty"`
}

type SearchForTransactionsParams struct {
	/**
	 * Only include transactions with this address in one of the transaction fields.
	 */
	Address string `url:"address,omitempty"`

	/**
	 * Combine with the address parameter to define what type of address to search for.
	 */
	AddressRole string `url:"address-role,omitempty"`

	/**
	 * Include results after the given time. Must be an RFC 3339 formatted string.
	 */
	AfterTime string `url:"after-time,omitempty"`

	/**
	 * Application ID
	 */
	ApplicationId uint64 `url:"application-id,omitempty"`

	/**
	 * Asset ID
	 */
	AssetId uint64 `url:"asset-id,omitempty"`

	/**
	 * Include results before the given time. Must be an RFC 3339 formatted string.
	 */
	BeforeTime string `url:"before-time,omitempty"`

	/**
	 * Results should have an amount greater than this value. MicroAlgos are the
	 * default currency unless an asset-id is provided, in which case the asset will be
	 * used.
	 */
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	/**
	 * Results should have an amount less than this value. MicroAlgos are the default
	 * currency unless an asset-id is provided, in which case the asset will be used.
	 */
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	/**
	 * Combine with address and address-role parameters to define what type of address
	 * to search for. The close to fields are normally treated as a receiver, if you
	 * would like to exclude them set this parameter to true.
	 */
	ExcludeCloseTo bool `url:"exclude-close-to,omitempty"`

	/**
	 * Maximum number of results to return.
	 */
	Limit uint64 `url:"limit,omitempty"`

	/**
	 * Include results at or before the specified max-round.
	 */
	MaxRound uint64 `url:"max-round,omitempty"`

	/**
	 * Include results at or after the specified min-round.
	 */
	MinRound uint64 `url:"min-round,omitempty"`

	/**
	 * The next page of results. Use the next token provided by the previous results.
	 */
	Next string `url:"next,omitempty"`

	/**
	 * Specifies a prefix which must be contained in the note field.
	 */
	NotePrefix string `url:"note-prefix,omitempty"`

	/**
	 * Include results which include the rekey-to field.
	 */
	RekeyTo bool `url:"rekey-to,omitempty"`

	/**
	 * Include results for the specified round.
	 */
	Round uint64 `url:"round,omitempty"`

	/**
	 * SigType filters just results using the specified type of signature:
	 *   sig - Standard
	 *   msig - MultiSig
	 *   lsig - LogicSig
	 */
	SigType string `url:"sig-type,omitempty"`

	TxType string `url:"tx-type,omitempty"`

	/**
	 * Lookup the specific transaction by ID.
	 */
	Txid string `url:"txid,omitempty"`
}
