//go:build falcon

package test

import (
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func genFalconKey() error {
	falconSigner = crypto.GenerateFalcon1024Account().AsSigner()
	return nil
}

func loadFalconKey() error {
	seed := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	// Address: AZM6UV2ONIVHH7BK2CSBUPJCXNPZH5LFA2YFBCZPHSYXUFJ4LLLFJOUT5Y
	falconAccount, err := crypto.Falcon1024AccountFromPQSeed(seed[:])
	if err != nil {
		return err
	}
	falconSigner = falconAccount.AsSigner()
	return err
}

func mnForFalcon(mn string) error {
	seed, err := mnemonic.ToPQSeed(mn, types.PQSchemeFalcon1024)
	if err != nil {
		return err
	}
	falconAccount, err := crypto.Falcon1024AccountFromPQSeed(seed)
	if err != nil {
		return err
	}
	falconSigner = falconAccount.AsSigner()
	return err
}
