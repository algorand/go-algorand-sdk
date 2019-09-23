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
// KeyReg parameters:
// - votePK is a base64-encoded string corresponding to the root participation public key
// - selectionKey is a base64-encoded string corresponding to the vrf public key
// - voteFirst is the first round this participation key is valid
// - voteLast is the last round this participation key is valid
// - voteKeyDilution is the dilution for the 2-level participation key
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
