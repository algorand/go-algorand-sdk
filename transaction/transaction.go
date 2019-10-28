package transaction

import (
	"encoding/base64"
	"fmt"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

const minFee = 1000

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
	total uint64, defaultFrozen bool, manager string, reserve string, freeze string, clawback string, unitName string, assetName string) (encoded []byte, err error) {
	var tx types.Transaction

	tx.Type = types.AssetConfigTx
	tx.AssetParams = types.AssetParams{
		Total:         total,
		DefaultFrozen: defaultFrozen,
	}

	if manager != "" {
		tx.AssetParams.Manager, err = types.DecodeAddress(manager)
		if err != nil {
			return
		}
	}
	if reserve != "" {
		tx.AssetParams.Reserve, err = types.DecodeAddress(reserve)
		if err != nil {
			return
		}
	}
	if freeze != "" {
		tx.AssetParams.Freeze, err = types.DecodeAddress(freeze)
		if err != nil {
			return
		}
	}
	if clawback != "" {
		tx.AssetParams.Clawback, err = types.DecodeAddress(clawback)
		if err != nil {
			return
		}
	}
	if len(unitName) > len(tx.AssetParams.UnitName) {
		err = fmt.Errorf("asset unit name %s too long (max %d bytes)", unitName, len(tx.AssetParams.UnitName))
		return
	}
	copy(tx.AssetParams.UnitName[:], []byte(unitName))

	if len(assetName) > len(tx.AssetParams.AssetName) {
		err = fmt.Errorf("asset name %s too long (max %d bytes)", assetName, len(tx.AssetParams.AssetName))
		return
	}
	copy(tx.AssetParams.AssetName[:], []byte(assetName))

	// Fill in header
	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return
	}
	ghBytes, err := byte32FromBase64(genesisHash)
	if err != nil {
		return
	}
	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         types.Algos(feePerByte),
		FirstValid:  types.Round(firstRound),
		LastValid:   types.Round(lastRound),
		GenesisHash: types.Digest(ghBytes),
		GenesisID:   genesisID,
		Note:        note,
	}
	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return
	}
	tx.Fee = types.Algos(eSize * feePerByte)

	if tx.Fee < minFee {
		tx.Fee = minFee
	}

	encoded = msgpack.Encode(tx)

	return
}

// MakeAssetConfigTxn creates a tx template for changing the
// keys for an asset. An empty string means a zero key (which
// cannot be changed after becoming zero); to keep a key
// unchanged, you must specify that key.
// - account is a checksummed, human-readable address for which we register the given participation key.
// - fee is a flat fee
// - firstRound is the first round this txn is valid (txn semantics unrelated to key registration)
// - lastRound is the last round this txn is valid
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
func MakeAssetConfigTxn(account string, feePerByte, firstRound, lastRound uint64, note []byte, genesisID string, genesisHash string, creator string,
	index uint64, newManager, newReserve, newFreeze, newClawback string) (encoded []byte, err error) {
	var tx types.Transaction

	tx.Type = types.AssetConfigTx

	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return
	}

	ghBytes, err := byte32FromBase64(genesisHash)
	if err != nil {
		return
	}

	tx.Header = types.Header{
		Sender:      accountAddr,
		Fee:         types.Algos(feePerByte),
		FirstValid:  types.Round(firstRound),
		LastValid:   types.Round(lastRound),
		GenesisHash: types.Digest(ghBytes),
		GenesisID:   genesisID,
	}

	creatorAddr, err := types.DecodeAddress(creator)
	if err != nil {
		return
	}

	tx.ConfigAsset = types.AssetID{
		Creator: creatorAddr,
		Index:   index,
	}

	if newManager != "" {
		tx.Type = types.AssetConfigTx
		tx.AssetParams.Manager, err = types.DecodeAddress(newManager)
		if err != nil {
			return
		}
	}

	if newReserve != "" {
		tx.AssetParams.Reserve, err = types.DecodeAddress(newReserve)
		if err != nil {
			return
		}
	}

	if newFreeze != "" {
		tx.AssetParams.Freeze, err = types.DecodeAddress(newFreeze)
		if err != nil {
			return
		}
	}

	if newClawback != "" {
		tx.AssetParams.Clawback, err = types.DecodeAddress(newClawback)
		if err != nil {
			return
		}
	}

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return
	}
	tx.Fee = types.Algos(eSize * feePerByte)

	if tx.Fee < minFee {
		tx.Fee = minFee
	}

	encoded = msgpack.Encode(tx)

	return
}

// MakePaymentTxn constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
func MakePaymentTxn(from, to string, fee, amount, firstRound, lastRound uint64, note []byte, closeRemainderTo, genesisID string, genesisHash []byte) (encoded []byte, err error) {
	// Decode from address
	fromAddr, err := types.DecodeAddress(from)
	if err != nil {
		return
	}

	// Decode to address
	toAddr, err := types.DecodeAddress(to)
	if err != nil {
		return
	}

	// Decode the CloseRemainderTo address, if present
	var closeRemainderToAddr types.Address
	if closeRemainderTo != "" {
		closeRemainderToAddr, err = types.DecodeAddress(closeRemainderTo)
		if err != nil {
			return
		}
	}

	// Decode GenesisHash
	if len(genesisHash) == 0 {
		err = fmt.Errorf("payment transaction must contain a genesisHash")
		return
	}

	var gh types.Digest
	copy(gh[:], genesisHash)

	// Build the transaction
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:      fromAddr,
			FirstValid:  types.Round(firstRound),
			LastValid:   types.Round(lastRound),
			Note:        note,
			GenesisID:   genesisID,
			GenesisHash: gh,
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver:         toAddr,
			Amount:           types.Algos(amount),
			CloseRemainderTo: closeRemainderToAddr,
		},
	}

	// Get the right fee
	l, err := estimateSize(tx)
	if err != nil {
		return
	}

	tx.Fee = types.Algos(uint64(fee) * l)

	if tx.Fee < minFee {
		tx.Fee = minFee
	}

	encoded = msgpack.Encode(tx)

	return
}

func estimateSize(tx types.Transaction) (uint64, error) {
	key := crypto.GenerateSK()
	en, err := crypto.SignTransaction(key, msgpack.Encode(tx))
	if err != nil {
		return 0, err
	}

	return uint64(len(en)), nil
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
