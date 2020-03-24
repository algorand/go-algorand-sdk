package future

import (
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// MinTxnFee is v5 consensus params, in microAlgos
const MinTxnFee = transaction.MinTxnFee

// MakePaymentTxn constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
// fee is fee per byte as received from algod SuggestedFee API call
func MakePaymentTxn(from, to string, amount uint64, note []byte, closeRemainderTo string, params types.SuggestedParams) (types.Transaction, error) {
	// Decode from address
	fromAddr, err := types.DecodeAddress(from)
	if err != nil {
		return types.Transaction{}, err
	}

	// Decode to address
	toAddr, err := types.DecodeAddress(to)
	if err != nil {
		return types.Transaction{}, err
	}

	// Decode the CloseRemainderTo address, if present
	var closeRemainderToAddr types.Address
	if closeRemainderTo != "" {
		closeRemainderToAddr, err = types.DecodeAddress(closeRemainderTo)
		if err != nil {
			return types.Transaction{}, err
		}
	}

	if len(params.GenesisHash) == 0 {
		return types.Transaction{}, fmt.Errorf("payment transaction must contain a genesisHash")
	}

	var gh types.Digest
	copy(gh[:], params.GenesisHash)

	// Build the transaction
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:      fromAddr,
			Fee:         params.Fee,
			FirstValid:  params.FirstRoundValid,
			LastValid:   params.LastRoundValid,
			Note:        note,
			GenesisID:   params.GenesisID,
			GenesisHash: gh,
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver:         toAddr,
			Amount:           types.MicroAlgos(amount),
			CloseRemainderTo: closeRemainderToAddr,
		},
	}

	// Update fee
	if !params.FlatFee {
		eSize, err := transaction.EstimateSize(tx)
		if err != nil {
			return types.Transaction{}, err
		}
		tx.Fee = types.MicroAlgos(eSize * uint64(params.Fee))
	}

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeKeyRegTxn constructs a keyreg transaction using the passed parameters.
// - account is a checksummed, human-readable address for which we register the given participation key.
// - note is a byte array
// - params is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
// KeyReg parameters:
// - votePK is a base64-encoded string corresponding to the root participation public key
// - selectionKey is a base64-encoded string corresponding to the vrf public key
// - voteFirst is the first round this participation key is valid
// - voteLast is the last round this participation key is valid
// - voteKeyDilution is the dilution for the 2-level participation key
func MakeKeyRegTxn(account string, note []byte, params types.SuggestedParams, voteKey, selectionKey string, voteFirst, voteLast, voteKeyDilution uint64) (types.Transaction, error) {
	// Decode account address
	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return types.Transaction{}, err
	}

	if len(params.GenesisHash) == 0 {
		return types.Transaction{}, fmt.Errorf("key registration transaction must contain a genesisHash")
	}

	var gh types.Digest
	copy(gh[:], params.GenesisHash)

	votePKBytes, err := byte32FromBase64(voteKey)
	if err != nil {
		return types.Transaction{}, err
	}

	selectionPKBytes, err := byte32FromBase64(selectionKey)
	if err != nil {
		return types.Transaction{}, err
	}

	tx := types.Transaction{
		Type: types.KeyRegistrationTx,
		Header: types.Header{
			Sender:      accountAddr,
			Fee:         params.Fee,
			FirstValid:  params.FirstRoundValid,
			LastValid:   params.LastRoundValid,
			Note:        note,
			GenesisHash: gh,
			GenesisID:   params.GenesisID,
		},
		KeyregTxnFields: types.KeyregTxnFields{
			VotePK:          types.VotePK(votePKBytes),
			SelectionPK:     types.VRFPK(selectionPKBytes),
			VoteFirst:       types.Round(voteFirst),
			VoteLast:        types.Round(voteLast),
			VoteKeyDilution: voteKeyDilution,
		},
	}

	if !params.FlatFee {
		// Update fee
		eSize, err := transaction.EstimateSize(tx)
		if err != nil {
			return types.Transaction{}, err
		}
		tx.Fee = types.MicroAlgos(eSize * uint64(params.Fee))
	}

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetCreateTxn constructs an asset creation transaction using the passed parameters.
// - account is a checksummed, human-readable address which will send the transaction.
// - note is a byte array
// - params is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
// Asset creation parameters:
// - see asset.go
func MakeAssetCreateTxn(account string, note []byte, params types.SuggestedParams, total uint64, decimals uint32, defaultFrozen bool, manager, reserve, freeze, clawback string, unitName, assetName, url, metadataHash string) (types.Transaction, error) {
	var tx types.Transaction
	var err error

	if decimals > types.AssetMaxNumberOfDecimals {
		return tx, fmt.Errorf("cannot create an asset with number of decimals %d (more than maximum %d)", decimals, types.AssetMaxNumberOfDecimals)
	}

	tx.Type = types.AssetConfigTx
	tx.AssetParams = types.AssetParams{
		Total:         total,
		Decimals:      decimals,
		DefaultFrozen: defaultFrozen,
		UnitName:      unitName,
		AssetName:     assetName,
		URL:           url,
	}

	if manager != "" {
		tx.AssetParams.Manager, err = types.DecodeAddress(manager)
		if err != nil {
			return tx, err
		}
	}
	if reserve != "" {
		tx.AssetParams.Reserve, err = types.DecodeAddress(reserve)
		if err != nil {
			return tx, err
		}
	}
	if freeze != "" {
		tx.AssetParams.Freeze, err = types.DecodeAddress(freeze)
		if err != nil {
			return tx, err
		}
	}
	if clawback != "" {
		tx.AssetParams.Clawback, err = types.DecodeAddress(clawback)
		if err != nil {
			return tx, err
		}
	}

	if len(assetName) > types.AssetNameMaxLen {
		return tx, fmt.Errorf("asset name too long: %d > %d", len(assetName), types.AssetNameMaxLen)
	}
	tx.AssetParams.AssetName = assetName

	if len(url) > types.AssetURLMaxLen {
		return tx, fmt.Errorf("asset url too long: %d > %d", len(url), types.AssetURLMaxLen)
	}
	tx.AssetParams.URL = url

	if len(unitName) > types.AssetUnitNameMaxLen {
		return tx, fmt.Errorf("asset unit name too long: %d > %d", len(unitName), types.AssetUnitNameMaxLen)
	}
	tx.AssetParams.UnitName = unitName

	if len(metadataHash) > types.AssetMetadataHashLen {
		return tx, fmt.Errorf("asset metadata hash '%s' too long: %d > %d)", metadataHash, len(metadataHash), types.AssetMetadataHashLen)
	}
	copy(tx.AssetParams.MetadataHash[:], []byte(metadataHash))

	if len(params.GenesisHash) == 0 {
		return types.Transaction{}, fmt.Errorf("asset transaction must contain a genesisHash")
	}
	var gh types.Digest
	copy(gh[:], params.GenesisHash)

	// Fill in header
	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         params.Fee,
		FirstValid:  params.FirstRoundValid,
		LastValid:   params.LastRoundValid,
		GenesisHash: gh,
		GenesisID:   params.GenesisID,
		Note:        note,
	}

	// Update fee
	if !params.FlatFee {
		eSize, err := transaction.EstimateSize(tx)
		if err != nil {
			return types.Transaction{}, err
		}
		tx.Fee = types.MicroAlgos(eSize * uint64(params.Fee))
	}

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetConfigTxn creates a tx template for changing the
// key configuration of an existing asset.
// Important notes -
// 	* Every asset config transaction is a fresh one. No parameters will be inherited from the current config.
// 	* Once an address is set to to the empty string, IT CAN NEVER BE CHANGED AGAIN. For example, if you want to keep
//    The current manager, you must specify its address again.
//	Parameters -
// - account is a checksummed, human-readable address that will send the transaction
// - note is an arbitrary byte array
// - params is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
// - index is the asset index id
// - for newManager, newReserve, newFreeze, newClawback see asset.go
// - strictEmptyAddressChecking: if true, disallow empty admin accounts from being set (preventing accidental disable of admin features)
func MakeAssetConfigTxn(account string, note []byte, params types.SuggestedParams, index uint64, newManager, newReserve, newFreeze, newClawback string, strictEmptyAddressChecking bool) (types.Transaction, error) {
	var tx types.Transaction

	if strictEmptyAddressChecking && (newManager == "" || newReserve == "" || newFreeze == "" || newClawback == "") {
		return tx, fmt.Errorf("strict empty address checking requested but empty address supplied to one or more manager addresses")
	}

	tx.Type = types.AssetConfigTx

	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return tx, err
	}

	if len(params.GenesisHash) == 0 {
		return types.Transaction{}, fmt.Errorf("asset transaction must contain a genesisHash")
	}
	var gh types.Digest
	copy(gh[:], params.GenesisHash)

	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         params.Fee,
		FirstValid:  params.FirstRoundValid,
		LastValid:   params.LastRoundValid,
		GenesisHash: gh,
		GenesisID:   params.GenesisID,
		Note:        note,
	}

	tx.ConfigAsset = types.AssetIndex(index)

	if newManager != "" {
		tx.Type = types.AssetConfigTx
		tx.AssetParams.Manager, err = types.DecodeAddress(newManager)
		if err != nil {
			return tx, err
		}
	}

	if newReserve != "" {
		tx.AssetParams.Reserve, err = types.DecodeAddress(newReserve)
		if err != nil {
			return tx, err
		}
	}

	if newFreeze != "" {
		tx.AssetParams.Freeze, err = types.DecodeAddress(newFreeze)
		if err != nil {
			return tx, err
		}
	}

	if newClawback != "" {
		tx.AssetParams.Clawback, err = types.DecodeAddress(newClawback)
		if err != nil {
			return tx, err
		}
	}

	if !params.FlatFee {
		// Update fee
		eSize, err := transaction.EstimateSize(tx)
		if err != nil {
			return types.Transaction{}, err
		}
		tx.Fee = types.MicroAlgos(eSize * uint64(params.Fee))
	}

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// transferAssetBuilder is a helper that builds asset transfer transactions:
// either a normal asset transfer, or an asset revocation
func transferAssetBuilder(account, recipient string, amount uint64, note []byte, params types.SuggestedParams, index uint64, closeAssetsTo, revocationTarget string) (types.Transaction, error) {
	var tx types.Transaction
	tx.Type = types.AssetTransferTx

	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return tx, err
	}

	if len(params.GenesisHash) == 0 {
		return types.Transaction{}, fmt.Errorf("asset transaction must contain a genesisHash")
	}
	var gh types.Digest
	copy(gh[:], params.GenesisHash)

	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         params.Fee,
		FirstValid:  params.FirstRoundValid,
		LastValid:   params.LastRoundValid,
		GenesisHash: gh,
		GenesisID:   params.GenesisID,
		Note:        note,
	}

	tx.XferAsset = types.AssetIndex(index)

	recipientAddr, err := types.DecodeAddress(recipient)
	if err != nil {
		return tx, err
	}
	tx.AssetReceiver = recipientAddr

	if closeAssetsTo != "" {
		closeToAddr, err := types.DecodeAddress(closeAssetsTo)
		if err != nil {
			return tx, err
		}
		tx.AssetCloseTo = closeToAddr
	}

	if revocationTarget != "" {
		revokedAddr, err := types.DecodeAddress(revocationTarget)
		if err != nil {
			return tx, err
		}
		tx.AssetSender = revokedAddr
	}

	tx.AssetAmount = amount

	// Update fee
	eSize, err := transaction.EstimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * uint64(params.Fee))

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetTransferTxn creates a tx for sending some asset from an asset holder to another user
// the recipient address must have previously issued an asset acceptance transaction for this asset
// - account is a checksummed, human-readable address that will send the transaction and assets
// - recipient is a checksummed, human-readable address what will receive the assets
// - amount is the number of assets to send
// - note is an arbitrary byte array
// - params is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
// - closeAssetsTo is a checksummed, human-readable address that behaves as a close-to address for the asset transaction; the remaining assets not sent to recipient will be sent to closeAssetsTo. Leave blank for no close-to behavior.
// - index is the asset index
func MakeAssetTransferTxn(account, recipient string, amount uint64, note []byte, params types.SuggestedParams, closeAssetsTo string, index uint64) (types.Transaction, error) {
	revocationTarget := "" // no asset revocation, this is normal asset transfer
	return transferAssetBuilder(account, recipient, amount, note, params, index, closeAssetsTo, revocationTarget)
}

// MakeAssetAcceptanceTxn creates a tx for marking an account as willing to accept the given asset
// - account is a checksummed, human-readable address that will send the transaction and begin accepting the asset
// - note is an arbitrary byte array
// - params is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
// - index is the asset index
func MakeAssetAcceptanceTxn(account string, note []byte, params types.SuggestedParams, index uint64) (types.Transaction, error) {
	return MakeAssetTransferTxn(account, account, 0, note, params, "", index)
}

// MakeAssetRevocationTxn creates a tx for revoking an asset from an account and sending it to another
// - account is a checksummed, human-readable address; it must be the revocation manager / clawback address from the asset's parameters
// - target is a checksummed, human-readable address; it is the account whose assets will be revoked
// - recipient is a checksummed, human-readable address; it will receive the revoked assets
// - amount defines the number of assets to clawback
// - params is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
// - index is the asset index
func MakeAssetRevocationTxn(account, target string, amount uint64, recipient string, note []byte, params types.SuggestedParams, index uint64) (types.Transaction, error) {
	closeAssetsTo := "" // no close-out, this is an asset revocation
	return transferAssetBuilder(account, recipient, amount, note, params, index, closeAssetsTo, target)
}

// MakeAssetDestroyTxn creates a tx template for destroying an asset, removing it from the record.
// All outstanding asset amount must be held by the creator, and this transaction must be issued by the asset manager.
// - account is a checksummed, human-readable address that will send the transaction; it also must be the asset manager
// - params is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
// - index is the asset index
func MakeAssetDestroyTxn(account string, note []byte, params types.SuggestedParams, index uint64) (types.Transaction, error) {
	// an asset destroy transaction is just a configuration transaction with AssetParams zeroed
	return MakeAssetConfigTxn(account, note, params, index, "", "", "", "", false)
}

// MakeAssetFreezeTxn constructs a transaction that freezes or unfreezes an account's asset holdings
// It must be issued by the freeze address for the asset
// - account is a checksummed, human-readable address which will send the transaction.
// - note is an optional arbitrary byte array
// - params is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
// - assetIndex is the index for tracking the asset
// - target is the account to be frozen or unfrozen
// - newFreezeSetting is the new state of the target account
func MakeAssetFreezeTxn(account string, note []byte, params types.SuggestedParams, assetIndex uint64, target string, newFreezeSetting bool) (types.Transaction, error) {
	var tx types.Transaction

	tx.Type = types.AssetFreezeTx

	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return tx, err
	}

	if len(params.GenesisHash) == 0 {
		return types.Transaction{}, fmt.Errorf("asset transaction must contain a genesisHash")
	}
	var gh types.Digest
	copy(gh[:], params.GenesisHash)

	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         params.Fee,
		FirstValid:  params.FirstRoundValid,
		LastValid:   params.LastRoundValid,
		GenesisHash: gh,
		GenesisID:   params.GenesisID,
		Note:        note,
	}

	tx.FreezeAsset = types.AssetIndex(assetIndex)

	tx.FreezeAccount, err = types.DecodeAddress(target)
	if err != nil {
		return tx, err
	}

	tx.AssetFrozen = newFreezeSetting

	if !params.FlatFee {
		// Update fee
		eSize, err := transaction.EstimateSize(tx)
		if err != nil {
			return types.Transaction{}, err
		}
		tx.Fee = types.MicroAlgos(eSize * uint64(params.Fee))
	}

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// byte32FromBase64 decodes the input base64 string and outputs a
// 32 byte array, erroring if the input is the wrong length.
func byte32FromBase64(in string) (out [32]byte, err error) {
	slice, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return
	}
	if len(slice) != 32 {
		return out, fmt.Errorf("Input is not 32 bytes")
	}
	copy(out[:], slice)
	return
}
