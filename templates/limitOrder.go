package templates

import (
	"encoding/base64"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

// LimitOrder represents a swap between Algos and an Asset at some ratio or better.
//
// Deprecated: Use TealCompile source compilation instead.
type LimitOrder struct {
	ContractTemplate
	assetID uint64
	owner   string
}

// GetSwapAssetsTransaction returns a group transaction array which transfer funds according to the contract's ratio
// assetAmount: amount of assets to be sent
// contract: byteform of the contract from the payer
// secretKey: secret key for signing transactions
// microAlgoAmount: number of microAlgos to transfer
// params: txn params for the transactions
// the first payment sends money (Algos) from contract to the recipient (we'll call him Buyer), closing the rest of the account to Owner
// the second payment sends money (the asset) from Buyer to the Owner
// these transactions will be rejected if they do not meet the restrictions set by the contract
//
// Deprecated: Use TealCompile source compilation instead.
func (lo LimitOrder) GetSwapAssetsTransaction(assetAmount uint64, microAlgoAmount uint64, contract, secretKey []byte, params types.SuggestedParams) ([]byte, error) {
	var buyerAddress types.Address
	copy(buyerAddress[:], secretKey[32:])
	contractAddress := crypto.AddressFromProgram(contract)
	algosForAssets, err := future.MakePaymentTxn(contractAddress.String(), buyerAddress.String(), microAlgoAmount, nil, "", params)
	if err != nil {
		return nil, err
	}
	assetsForAlgos, err := future.MakeAssetTransferTxn(buyerAddress.String(), lo.owner, assetAmount, nil, params, "", lo.assetID)
	if err != nil {
		return nil, err
	}

	gid, err := crypto.ComputeGroupID([]types.Transaction{algosForAssets, assetsForAlgos})
	if err != nil {
		return nil, err
	}
	algosForAssets.Group = gid
	assetsForAlgos.Group = gid

	logicSig, err := crypto.MakeLogicSig(contract, nil, nil, crypto.MultisigAccount{})
	if err != nil {
		return nil, err
	}
	_, algosForAssetsSigned, err := crypto.SignLogicsigTransaction(logicSig, algosForAssets)
	if err != nil {
		return nil, err
	}
	_, assetsForAlgosSigned, err := crypto.SignTransaction(secretKey, assetsForAlgos)
	if err != nil {
		return nil, err
	}

	var signedGroup []byte
	signedGroup = append(signedGroup, algosForAssetsSigned...)
	signedGroup = append(signedGroup, assetsForAlgosSigned...)

	return signedGroup, nil
}

// MakeLimitOrder allows a user to exchange some number of assets for some number of algos.
// Fund the contract with some number of Algos to limit the maximum number of
// Algos you're willing to trade for some other asset.
//
// Works on two cases:
// * trading Algos for some other asset
// * closing out Algos back to the originator after a timeout
//
// trade case, a 2 transaction group:
// gtxn[0] (this txn) Algos from Me to Other
// gtxn[1] asset from Other to Me
//
// We want to get _at least_ some amount of the other asset per our Algos
// gtxn[1].AssetAmount / gtxn[0].Amount >= N / D
// ===
// gtxn[1].AssetAmount * D >= gtxn[0].Amount * N
//
// close-out case:
// txn alone, close out value after timeout
//
// Parameters:
//  - owner: the address to refund funds to on timeout
//  - assetID: ID of the transferred asset
//  - ratn: exchange rate (N asset per D Algos, or better)
//  - ratd: exchange rate (N asset per D Algos, or better)
//  - expiryRound: the round at which the account expires
//  - minTrade: the minimum amount (of Algos) to be traded away
//  - maxFee: maximum fee used by the limit order transaction
//
// Deprecated: Use TealCompile source compilation instead.
func MakeLimitOrder(owner string, assetID, ratn, ratd, expiryRound, minTrade, maxFee uint64) (LimitOrder, error) {
	const referenceProgram = "ASAKAAEFAgYEBwgJCiYBIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITMRYiEjEQIxIQMQEkDhAyBCMSQABVMgQlEjEIIQQNEDEJMgMSEDMBECEFEhAzAREhBhIQMwEUKBIQMwETMgMSEDMBEiEHHTUCNQExCCEIHTUENQM0ATQDDUAAJDQBNAMSNAI0BA8QQAAWADEJKBIxAiEJDRAxBzIDEhAxCCISEBA="
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return LimitOrder{}, err
	}

	var referenceOffsets = []uint64{ /*maxFee*/ 5 /*minTrade*/, 7 /*assetID*/, 9 /*ratd*/, 10 /*ratn*/, 11 /*expiryRound*/, 12 /*ownerAddr*/, 16}
	ownerAddr, err := types.DecodeAddress(owner)
	if err != nil {
		return LimitOrder{}, err
	}
	injectionVector := []interface{}{maxFee, minTrade, assetID, ratd, ratn, expiryRound, ownerAddr}
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return LimitOrder{}, err
	}

	address := crypto.AddressFromProgram(injectedBytes)
	lo := LimitOrder{
		ContractTemplate: ContractTemplate{
			address: address.String(),
			program: injectedBytes,
		},
		owner:   owner,
		assetID: assetID,
	}
	return lo, err
}
