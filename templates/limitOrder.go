package templates

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

// Split template representation
type LimitOrder struct {
	ContractTemplate
	ratn uint64
	ratd uint64
}

// GetSwapAssetsTransaction returns a group transaction array which transfer funds according to the contract's ratio
// amount: amount of assets to be sent
// contract: byteform of the contract from the payer
// secretKey: secret key for signing transactions
func GetSwapAssetsTransaction() ([]byte, error) {

	var stx1 []byte
	var stx2 []byte
	var signedGroup []byte
	signedGroup = append(signedGroup, stx1...)
	signedGroup = append(signedGroup, stx2...)

	return signedGroup, nil
}

// MakeLimitOrder
// Works on a contract account.
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
func MakeLimitOrder(owner string, assetID, ratn, ratd, expiryRound, minTrade, maxFee uint64) (LimitOrder, error) {
	const referenceProgram = "ASAKAAEFAgYEBwgJCiYBIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITMRYiEjEQIxIQMQEkDhAyBCMSQABVMgQlEjEIIQQNEDEJMgMSEDMBECEFEhAzAREhBhIQMwEUKBIQMwETMgMSEDMBEiEHHTUCNQExCCEIHTUENQM0ATQDDUAAJDQBNAMSNAI0BA8QQAAWADEJKBIxAiEJDRAxBzIDEhAxCCISEBA="
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return LimitOrder{}, err
	}

	var referenceOffsets = []uint64{ /*fee*/ 4 /*timeout*/, 7 /*ratn*/, 8 /*ratd*/, 9 /*minPay*/, 10 /*owner*/, 14 /*receiver1*/, 47 /*receiver2*/, 80}
	ownerAddr, err := types.DecodeAddress(owner)
	if err != nil {
		return LimitOrder{}, err
	}
	injectionVector := []interface{}{maxFee, expiryRound, ratn, ratd, minTrade, ownerAddr}
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
		ratn: ratn,
		ratd: ratd,
	}
	return lo, err
}
