package transaction

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

// MinTxnFee is v5 consensus params, in microAlgos
const MinTxnFee = 1000

// MakePaymentTxn constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
// fee is fee per byte as received from algod SuggestedFee API call
func MakePaymentTxn(from, to string, fee, amount, firstRound, lastRound uint64, note []byte, closeRemainderTo, genesisID string, genesisHash []byte) (types.Transaction, error) {
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

	// Decode GenesisHash
	if len(genesisHash) == 0 {
		return types.Transaction{}, fmt.Errorf("payment transaction must contain a genesisHash")
	}

	var gh types.Digest
	copy(gh[:], genesisHash)

	// Build the transaction
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:      fromAddr,
			Fee:         types.MicroAlgos(fee),
			FirstValid:  types.Round(firstRound),
			LastValid:   types.Round(lastRound),
			Note:        note,
			GenesisID:   genesisID,
			GenesisHash: gh,
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver:         toAddr,
			Amount:           types.MicroAlgos(amount),
			CloseRemainderTo: closeRemainderToAddr,
		},
	}

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakePaymentTxnWithFlatFee constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
// fee is a flat fee
func MakePaymentTxnWithFlatFee(from, to string, fee, amount, firstRound, lastRound uint64, note []byte, closeRemainderTo, genesisID string, genesisHash []byte) (types.Transaction, error) {
	tx, err := MakePaymentTxn(from, to, fee, amount, firstRound, lastRound, note, closeRemainderTo, genesisID, genesisHash)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeKeyRegTxn constructs a keyreg transaction using the passed parameters.
// - account is a checksummed, human-readable address for which we register the given participation key.
// - fee is fee per byte as received from algod SuggestedFee API call.
// - firstRound is the first round this txn is valid (txn semantics unrelated to key registration)
// - lastRound is the last round this txn is valid
// - note is a byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// KeyReg parameters:
// - votePK is a base64-encoded string corresponding to the root participation public key
// - selectionKey is a base64-encoded string corresponding to the vrf public key
// - voteFirst is the first round this participation key is valid
// - voteLast is the last round this participation key is valid
// - voteKeyDilution is the dilution for the 2-level participation key
func MakeKeyRegTxn(account string, feePerByte, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string,
	voteKey, selectionKey string, voteFirst, voteLast, voteKeyDilution uint64) (types.Transaction, error) {
	// Decode account address
	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return types.Transaction{}, err
	}

	ghBytes, err := byte32FromBase64(genesisHash)
	if err != nil {
		return types.Transaction{}, err
	}

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
			Fee:         types.MicroAlgos(feePerByte),
			FirstValid:  types.Round(firstRound),
			LastValid:   types.Round(lastRound),
			Note:        note,
			GenesisHash: types.Digest(ghBytes),
			GenesisID:   genesisID,
		},
		KeyregTxnFields: types.KeyregTxnFields{
			VotePK:          types.VotePK(votePKBytes),
			SelectionPK:     types.VRFPK(selectionPKBytes),
			VoteFirst:       types.Round(voteFirst),
			VoteLast:        types.Round(voteLast),
			VoteKeyDilution: voteKeyDilution,
		},
	}

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * feePerByte)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeKeyRegTxnWithFlatFee constructs a keyreg transaction using the passed parameters.
// - account is a checksummed, human-readable address for which we register the given participation key.
// - fee is a flat fee
// - firstRound is the first round this txn is valid (txn semantics unrelated to key registration)
// - lastRound is the last round this txn is valid
// - note is a byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// KeyReg parameters:
// - votePK is a base64-encoded string corresponding to the root participation public key
// - selectionKey is a base64-encoded string corresponding to the vrf public key
// - voteFirst is the first round this participation key is valid
// - voteLast is the last round this participation key is valid
// - voteKeyDilution is the dilution for the 2-level participation key
func MakeKeyRegTxnWithFlatFee(account string, fee, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string,
	voteKey, selectionKey string, voteFirst, voteLast, voteKeyDilution uint64) (types.Transaction, error) {
	tx, err := MakeKeyRegTxn(account, fee, firstRound, lastRound, note, genesisID, genesisHash, voteKey, selectionKey, voteFirst, voteLast, voteKeyDilution)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Fee = types.MicroAlgos(fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetCreateTxn constructs an asset creation transaction using the passed parameters.
// - account is a checksummed, human-readable address which will send the transaction.
// - fee is fee per byte as received from algod SuggestedFee API call.
// - firstRound is the first round this txn is valid (txn semantics unrelated to the asset)
// - lastRound is the last round this txn is valid
// - note is a byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// Asset creation parameters:
// - see asset.go
func MakeAssetCreateTxn(account string, feePerByte, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string,
	total uint64, defaultFrozen bool, manager string, reserve string, freeze string, clawback string, unitName string, assetName string) (types.Transaction, error) {
	var tx types.Transaction
	var err error

	tx.Type = types.AssetConfigTx
	tx.AssetParams = types.AssetParams{
		Total:         total,
		DefaultFrozen: defaultFrozen,
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
	if len(unitName) > len(tx.AssetParams.UnitName) {
		return tx, fmt.Errorf("asset unit name %s too long (max %d bytes)", unitName, len(tx.AssetParams.UnitName))
	}
	copy(tx.AssetParams.UnitName[:], []byte(unitName))

	if len(assetName) > len(tx.AssetParams.AssetName) {
		return tx, fmt.Errorf("asset name %s too long (max %d bytes)", assetName, len(tx.AssetParams.AssetName))
	}
	copy(tx.AssetParams.AssetName[:], []byte(assetName))

	// Fill in header
	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return types.Transaction{}, err
	}
	ghBytes, err := byte32FromBase64(genesisHash)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         types.MicroAlgos(feePerByte),
		FirstValid:  types.Round(firstRound),
		LastValid:   types.Round(lastRound),
		GenesisHash: types.Digest(ghBytes),
		GenesisID:   genesisID,
		Note:        note,
	}

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * feePerByte)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetConfigTxn creates a tx template for changing the
// keys for an asset. An empty string means a zero key (which
// cannot be changed after becoming zero); to keep a key
// unchanged, you must specify that key.
// - account is a checksummed, human-readable address that will send the transaction
// - fee is a fee per byte
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset config)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
func MakeAssetConfigTxn(account string, feePerByte, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string, creator string,
	index uint64, newManager, newReserve, newFreeze, newClawback string) (types.Transaction, error) {
	var tx types.Transaction

	tx.Type = types.AssetConfigTx

	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return tx, err
	}

	ghBytes, err := byte32FromBase64(genesisHash)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         types.MicroAlgos(feePerByte),
		FirstValid:  types.Round(firstRound),
		LastValid:   types.Round(lastRound),
		GenesisHash: types.Digest(ghBytes),
		GenesisID:   genesisID,
		Note:        note,
	}

	creatorAddr, err := types.DecodeAddress(creator)
	if err != nil {
		return tx, err
	}

	tx.ConfigAsset = types.AssetID{
		Creator: creatorAddr,
		Index:   index,
	}

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

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * feePerByte)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetTransferTxn creates a tx for sending some asset from an asset holder to another user
// the recipient address must have previously issued an asset acceptance transaction for this asset
// - account is a checksummed, human-readable address that will send the transaction and assets
// - recipient is a checksummed, human-readable address what will receive the assets
// - closeAssetsTo is a checksummed, human-readable address that behaves as a close-to address for the asset transaction; the remaining assets not sent to recipient will be sent to closeAssetsTo. Leave blank for no close-to behavior.
// - amount is the number of assets to send
// - feePerByte is a fee per byte
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetTransferTxn(account, recipient, closeAssetsTo string, amount, feePerByte, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	var tx types.Transaction

	tx.Type = types.AssetTransferTx

	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return tx, err
	}

	ghBytes, err := byte32FromBase64(genesisHash)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         types.MicroAlgos(feePerByte),
		FirstValid:  types.Round(firstRound),
		LastValid:   types.Round(lastRound),
		GenesisHash: types.Digest(ghBytes),
		GenesisID:   genesisID,
		Note:        note,
	}

	creatorAddr, err := types.DecodeAddress(creator)
	if err != nil {
		return tx, err
	}

	tx.XferAsset = types.AssetID{
		Creator: creatorAddr,
		Index:   index,
	}

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

	tx.AssetAmount = amount

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * feePerByte)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetAcceptanceTransaction creates a tx for marking an account as willing to accept the given asset
// - account is a checksummed, human-readable address that will send the transaction and begin accepting the asset
// - feePerByte is a fee per byte
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetAcceptanceTransaction(account string, feePerByte, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetTransferTxn(account, account, "", 0,
		feePerByte, firstRound, lastRound, note, genesisID, genesisHash, creator, index)
	return tx, err
}

// MakeAssetRevocationTransaction creates a tx for revoking an asset from an account and sending it to another
// - account is a checksummed, human-readable address; it must be the revocation manager / clawback address from the asset's parameters
// - target is a checksummed, human-readable address; it is the account whose assets will be revoked
// - recipient is a checksummed, human-readable address; it will receive the revoked assets
// - feePerByte is a fee per byte
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetRevocationTransaction(account, target, recipient string, amount, feePerByte, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetTransferTxn(account, recipient, "", amount,
		feePerByte, firstRound, lastRound, note, genesisID, genesisHash, creator, index)

	targetAddr, err := types.DecodeAddress(target)
	if err != nil {
		return tx, err
	}
	tx.AssetSender = targetAddr

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Fee = types.MicroAlgos(eSize * feePerByte)
	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, err
}

// MakeAssetMintTransaction creates a tx for minting more of an asset (transferring from the reserve account)
// Note that this is just an asset transfer transaction where sender=reserve
// - reserve is a checksummed, human-readable address that will send the transaction. it must be the asset's reserve to be considered a minting, but this is not enforced by the function.
// - recipient is a checksummed, human-readable address; it will receive the revoked assets
// - amount is the number of assets to mint
// - feePerByte is a fee per byte
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetMintTransaction(reserve, recipient string, amount, feePerByte, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetTransferTxn(reserve, recipient, "", amount,
		feePerByte, firstRound, lastRound, note, genesisID, genesisHash, creator, index)
	return tx, err
}

// MakeAssetBurnTransaction creates for burning an asset (transferring to the reserve account)
// Note that this is just an asset transfer transaction where recipient=reserve
// - account is a checksummed, human-readable address that will send the transaction.
// - reserve is a checksummed, human-readable address; it will receive the revoked assets and must be the reserve account to be considered a burn. this is not enforced by the function
// - amount is the number of assets to "burn"
// - feePerByte is a fee per byte
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetBurnTransaction(account, reserve string, amount, feePerByte, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetTransferTxn(account, reserve, "", amount,
		feePerByte, firstRound, lastRound, note, genesisID, genesisHash, creator, index)
	return tx, err
}

// MakeAssetDestroyTxn creates a tx template for destroying an asset, removing it from the record.
// All outstanding asset amount must be held by the creator, and this transaction must be issued by the asset manager.
// - account is a checksummed, human-readable address that will send the transaction; it also must be the asset manager
// - fee is a fee per byte
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetDestroyTxn(account string, feePerByte, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string,
	creator string, index uint64) (types.Transaction, error) {
	// an asset destroy transaction is just a configuration transaction with AssetParams zeroed
	tx, err := MakeAssetConfigTxn(account, feePerByte, firstRound, lastRound, note, genesisID, genesisHash,
		creator, index, "", "", "", "")

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * feePerByte)
	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetFreezeTxn constructs a transaction that freezes or unfreezes an account's asset holdings
// It must be issued by the freeze address for the asset
// - account is a checksummed, human-readable address which will send the transaction.
// - fee is fee per byte as received from algod SuggestedFee API call.
// - firstRound is the first round this txn is valid (txn semantics unrelated to the asset)
// - lastRound is the last round this txn is valid
// - note is an optional arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the creator address for the asset
// - assetIndex is the index for tracking the asset
// - target is the account to be frozen or unfrozen
// - newFreezeSetting is the new state of the target account
func MakeAssetFreezeTxn(account string, fee, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string,
	creator string, assetIndex uint64, target string, newFreezeSetting bool) (types.Transaction, error) {
	var tx types.Transaction

	tx.Type = types.AssetFreezeTx

	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return tx, err
	}

	ghBytes, err := byte32FromBase64(genesisHash)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         types.MicroAlgos(fee),
		FirstValid:  types.Round(firstRound),
		LastValid:   types.Round(lastRound),
		GenesisHash: types.Digest(ghBytes),
		GenesisID:   genesisID,
		Note:        note,
	}

	tx.FreezeAsset.Index = assetIndex
	tx.FreezeAsset.Creator, err = types.DecodeAddress(creator)
	if err != nil {
		return tx, err
	}

	tx.FreezeAccount, err = types.DecodeAddress(target)
	if err != nil {
		return tx, err
	}

	tx.AssetFrozen = newFreezeSetting
	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetCreateTxnWithFlatFee constructs an asset creation transaction using the passed parameters.
// - account is a checksummed, human-readable address which will send the transaction.
// - fee is fee per byte as received from algod SuggestedFee API call.
// - firstRound is the first round this txn is valid (txn semantics unrelated to the asset)
// - lastRound is the last round this txn is valid
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// Asset creation parameters:
// - see asset.go
func MakeAssetCreateTxnWithFlatFee(account string, fee, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string,
	total uint64, defaultFrozen bool, manager string, reserve string, freeze string, clawback string, unitName string, assetName string) (types.Transaction, error) {
	tx, err := MakeAssetCreateTxn(account, fee, firstRound, lastRound, note, genesisID, genesisHash, total, defaultFrozen, manager, reserve, freeze, clawback, unitName, assetName)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Fee = types.MicroAlgos(fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeAssetConfigTxnWithFlatFee creates a tx template for changing the
// keys for an asset. An empty string means a zero key (which
// cannot be changed after becoming zero); to keep a key
// unchanged, you must specify that key.
func MakeAssetConfigTxnWithFlatFee(account string, fee, firstRound, lastRound uint64, note []byte, genesisID, genesisHash, creator string,
	index uint64, newManager, newReserve, newFreeze, newClawback string) (types.Transaction, error) {
	tx, err := MakeAssetConfigTxn(account, fee, firstRound, lastRound, note, genesisID, genesisHash, creator, index, newManager, newReserve, newFreeze, newClawback)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Fee = types.MicroAlgos(fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}
	return tx, nil
}

// MakeAssetTransferTxnWithFlatFee creates a tx for sending some asset from an asset holder to another user
// the recipient address must have previously issued an asset acceptance transaction for this asset
// - account is a checksummed, human-readable address that will send the transaction and assets
// - recipient is a checksummed, human-readable address what will receive the assets
// - closeAssetsTo is a checksummed, human-readable address that behaves as a close-to address for the asset transaction; the remaining assets not sent to recipient will be sent to closeAssetsTo. Leave blank for no close-to behavior.
// - amount is the number of assets to send
// - fee is a flat fee
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetTransferTxnWithFlatFee(account, recipient, closeAssetsTo string, amount, fee, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetTransferTxn(account, recipient, closeAssetsTo, amount, fee, firstRound, lastRound, note, genesisID, genesisHash, creator, index)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Fee = types.MicroAlgos(fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}
	return tx, nil
}

// MakeAssetAcceptanceTransactionWithFlatFee creates a tx for marking an account as willing to accept an asset
// - account is a checksummed, human-readable address that will send the transaction and begin accepting the asset
// - fee is a flat fee
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetAcceptanceTransactionWithFlatFee(account string, fee, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetAcceptanceTransaction(account, fee, firstRound, lastRound, note, genesisID, genesisHash, creator, index)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Fee = types.MicroAlgos(fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}
	return tx, nil
}

// MakeAssetRevocationTransactionWithFlatFee creates a tx for revoking an asset from an account and sending it to another
// - account is a checksummed, human-readable address; it must be the revocation manager / clawback address from the asset's parameters
// - target is a checksummed, human-readable address; it is the account whose assets will be revoked
// - recipient is a checksummed, human-readable address; it will receive the revoked assets
// - fee is a flat fee
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetRevocationTransactionWithFlatFee(account, target, recipient string, amount, fee, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetRevocationTransaction(account, target, recipient, amount, fee, firstRound, lastRound,
		note, genesisID, genesisHash, creator, index)

	if err != nil {
		return types.Transaction{}, err
	}

	tx.Fee = types.MicroAlgos(fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}
	return tx, nil
}

// MakeAssetMintTransactionWithFlatFee creates a tx for minting more of an asset (transferring from the reserve account)
// Note that this is just an asset transfer transaction where sender=reserve
// - reserve is a checksummed, human-readable address that will send the transaction. it must be the asset's reserve to be considered a minting, but this is not enforced by the function.
// - recipient is a checksummed, human-readable address; it will receive the revoked assets
// - amount is the number of assets to send
// - fee is the flat fee
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetMintTransactionWithFlatFee(reserve, recipient string, amount, fee, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetTransferTxnWithFlatFee(reserve, recipient, "", amount,
		fee, firstRound, lastRound, note, genesisID, genesisHash, creator, index)
	return tx, err
}

// MakeAssetBurnTransactionWithFlatFee creates for burning an asset (transferring to the reserve account)
// Note that this is just an asset transfer transaction where recipient=reserve
// - account is a checksummed, human-readable address that will send the transaction.
// - reserve is a checksummed, human-readable address; it will receive the revoked assets and must be the reserve account to be considered a burn. this is not enforced by the function
// - amount is the number of assets to send
// - fee is the flat fee
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - note is an arbitrary byte array
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetBurnTransactionWithFlatFee(account, reserve string, amount, fee, firstRound, lastRound uint64, note []byte,
	genesisID, genesisHash, creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetTransferTxnWithFlatFee(account, reserve, "", amount,
		fee, firstRound, lastRound, note, genesisID, genesisHash, creator, index)
	return tx, err
}

// MakeAssetDestroyTxnWithFlatFee creates a tx template for destroying an asset, removing it from the record.
// All outstanding asset amount must be held by the creator, and this transaction must be issued by the asset manager.
// - account is a checksummed, human-readable address that will send the transaction; it also must be the asset manager
// - fee is a flat fee
// - firstRound is the first round this txn is valid (txn semantics unrelated to asset management)
// - lastRound is the last round this txn is valid
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// - index is the asset index
func MakeAssetDestroyTxnWithFlatFee(account string, fee, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string,
	creator string, index uint64) (types.Transaction, error) {
	tx, err := MakeAssetConfigTxnWithFlatFee(account, fee, firstRound, lastRound, note, genesisID, genesisHash,
		creator, index, "", "", "", "")
	return tx, err
}

// MakeAssetFreezeTxnWithFlatFee is as MakeAssetFreezeTxn, but taking a flat fee rather than a fee per byte.
func MakeAssetFreezeTxnWithFlatFee(account string, fee, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string,
	creator string, assetIndex uint64, target string, newFreezeSetting bool) (types.Transaction, error) {
	tx, err := MakeAssetFreezeTxn(account, fee, firstRound, lastRound, note, genesisID, genesisHash, creator, assetIndex, target, newFreezeSetting)
	if err != nil {
		return types.Transaction{}, err
	}

	tx.Fee = types.MicroAlgos(fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}
	return tx, nil
}

// AssignGroupID computes and return list of transactions with Group field set.
// - txns is a list of transactions to process
// - account specifies a sender field of transaction to return. Set to empty string to return all of them
func AssignGroupID(txns []types.Transaction, account string) (result []types.Transaction, err error) {
	gid, err := crypto.ComputeGroupID(txns)
	if err != nil {
		return
	}
	var decoded types.Address
	if account != "" {
		decoded, err = types.DecodeAddress(account)
		if err != nil {
			return
		}
	}
	for _, tx := range txns {
		if account == "" || bytes.Compare(tx.Sender[:], decoded[:]) == 0 {
			tx.Group = gid
			result = append(result, tx)
		}
	}
	return result, nil
}

// EstimateSize returns the estimated length of the encoded transaction
func estimateSize(txn types.Transaction) (uint64, error) {
	key := crypto.GenerateAccount()
	_, stx, err := crypto.SignTransaction(key.PrivateKey, txn)
	if err != nil {
		return 0, err
	}
	return uint64(len(stx)), nil
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
