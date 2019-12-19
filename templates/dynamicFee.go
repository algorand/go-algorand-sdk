package templates

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

// DynamicFee template representation
type DynamicFee struct {
	ContractTemplate
}

// MakeDynamicFee contract allows you to create a transaction without
// specifying the fee. The fee will be determined at the moment of
// transfer.
//
// Parameters:
//  - receiver: address which is authorized to receive withdrawals
//  - amount: the maximum number of funds allowed for a single withdrawal
//  - withdrawWindow: the duration of a withdrawal period
//  - period: the time between a pair of withdrawal periods
//  - expiryRound: the round at which the account expires
//  - maxFee: maximum fee used by the withdrawal transaction
func MakeDynamicFee(receiver, closeRemainder string, amount, firstValid, lastValid uint64) (DynamicFee, error) {
	leaseBytes := make([]byte, 32)
	crypto.RandomBytes(leaseBytes)
	leaseString := base64.StdEncoding.EncodeToString(leaseBytes)
	return MakeDynamicFeeWithLease(receiver, closeRemainder, leaseString, amount, firstValid, lastValid)
}

// MakeDynamicFeeWithLease is as MakeDynamicFee, but the caller can specify the lease (using b64 string)
func MakeDynamicFeeWithLease(receiver, closeRemainder, lease string, amount, firstValid, lastValid uint64) (DynamicFee, error) {
	const referenceProgram = "ASAFAgEHBgUmAyD+vKC7FEpaTqe0OKRoGsgObKEFvLYH/FZTJclWlfaiEyDmmpYeby1feshmB5JlUr6YI17TM2PKiJGLuck4qRW2+SB/g7Flf/H8U7ktwYFIodZd/C1LH6PWdyhK3dIAEm2QaTIEIhIzABAjEhAzAAcxABIQMwAIMQESEDEWIxIQMRAjEhAxBygSEDEJKRIQMQgkEhAxAiUSEDEEIQQSEDEGKhIQ"
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return DynamicFee{}, err
	}
	receiverAddr, err := types.DecodeAddress(receiver)
	if err != nil {
		return DynamicFee{}, err
	}
	closeRemainderAddr, err := types.DecodeAddress(closeRemainder)
	if err != nil {
		return DynamicFee{}, err
	}

	var referenceOffsets = []uint64{ /*amount*/ 5 /*firstValid*/, 6 /*lastValid*/, 7 /*receiver*/, 11 /*closeRemainder*/, 44 /*lease*/, 76}
	injectionVector := []interface{}{amount, firstValid, lastValid, receiverAddr, closeRemainderAddr, lease}
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return DynamicFee{}, err
	}

	address := crypto.AddressFromProgram(injectedBytes)
	dynamicFee := DynamicFee{
		ContractTemplate: ContractTemplate{
			address: address.String(),
			program: injectedBytes,
		},
	}
	return dynamicFee, err
}
