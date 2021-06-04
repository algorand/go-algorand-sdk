package test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cucumber/godog"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
)

var sk1 ed25519.PrivateKey
var addr1 types.Address

func aSigningAccountWithAddressAndMnemonic(address, mnem string) error {
	var err error
	addr1, err = types.DecodeAddress(address)
	if err != nil {
		return err
	}
	sk1, err = mnemonic.ToPrivateKey(mnem)
	return err
}

func iBuildAnApplicationTransactionUnit(
	operation string,
	applicationIdInt int,
	sender, approvalProgram, clearProgram string,
	globalBytes, globalInts, localBytes, localInts int,
	appArgs, foreignApps, foreignAssets, appAccounts string,
	extraPages int,
	fee, firstValid, lastValid int,
	genesisHash string) error {

	applicationId = uint64(applicationIdInt)
	var clearP []byte
	var approvalP []byte
	var err error

	if approvalProgram != "" {
		approvalP, err = ioutil.ReadFile("features/resources/" + approvalProgram)
		if err != nil {
			return err
		}
	}

	if clearProgram != "" {
		clearP, err = ioutil.ReadFile("features/resources/" + clearProgram)
		if err != nil {
			return err
		}
	}
	args, err := parseAppArgs(appArgs)
	if err != nil {
		return err
	}
	var accs []string
	if appAccounts != "" {
		accs = strings.Split(appAccounts, ",")
	}

	fApp, err := splitUint64(foreignApps)
	if err != nil {
		return err
	}

	fAssets, err := splitUint64(foreignAssets)
	if err != nil {
		return err
	}

	gSchema := types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}
	lSchema := types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)}

	suggestedParams, err := getSuggestedParams(uint64(fee), uint64(firstValid), uint64(lastValid), "", genesisHash, true)
	if err != nil {
		return err
	}

	switch operation {
	case "create":
		tx, err = future.MakeApplicationCreateTx(false, approvalP, clearP,
			gSchema, lSchema, args, accs, fApp, fAssets, uint32(extraPages),
			suggestedParams, addr1, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "update":
		tx, err = future.MakeApplicationUpdateTx(applicationId, args, accs, fApp, fAssets,
			approvalP, clearP,
			suggestedParams, addr1, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "call":
		tx, err = future.MakeApplicationCallTx(applicationId, args, accs,
			fApp, fAssets, types.NoOpOC, approvalP, clearP, gSchema, lSchema, 0,
			suggestedParams, addr1, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "optin":
		tx, err = future.MakeApplicationOptInTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, addr1, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "clear":
		tx, err = future.MakeApplicationClearStateTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, addr1, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "closeout":
		tx, err = future.MakeApplicationCloseOutTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, addr1, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "delete":
		tx, err = future.MakeApplicationDeleteTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, addr1, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}
	}
	return nil

}

func signTheTransaction() error {
	var err error
	_, stx, err = crypto.SignTransaction(sk1, tx)
	return err
}

func theBaseEncodedSignedTransactionShouldEqual(base int, golden string) error {
	gold, err := base64.StdEncoding.DecodeString(golden)
	if err != nil {
		return err
	}
	if !bytes.Equal(gold, stx) {
		return fmt.Errorf("Application signed transaction does not match the golden.")
	}
	return nil
}

func getSuggestedParams(
	fee, fv, lv uint64,
	gen, ghb64 string,
	flat bool) (types.SuggestedParams, error) {
	gh, err := base64.StdEncoding.DecodeString(ghb64)
	if err != nil {
		return types.SuggestedParams{}, err
	}
	return types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(fv),
		LastRoundValid:  types.Round(lv),
		FlatFee:         flat,
	}, err
}

func weMakeAGetAssetByIDCall(assetID int) error {
	clt, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	clt.GetAssetByID(uint64(assetID)).Do(context.Background())
	return nil
}

func weMakeAGetApplicationByIDCall(applicationID int) error {
	clt, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	clt.GetApplicationByID(uint64(applicationID)).Do(context.Background())
	return nil
}

func weMakeASearchForApplicationsCall(applicationID int) error {
	clt, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	clt.SearchForApplications().ApplicationId(uint64(applicationID)).Do(context.Background())
	return nil
}

func weMakeALookupApplicationsCall(applicationID int) error {
	clt, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	clt.LookupApplicationByID(uint64(applicationID)).Do(context.Background())
	return nil
}

func ApplicationsUnitContext(s *godog.Suite) {
	// @unit.transactiosn
	s.Step(`^a signing account with address "([^"]*)" and mnemonic "([^"]*)"$`, aSigningAccountWithAddressAndMnemonic)
	s.Step(`^I build an application transaction with operation "([^"]*)", application-id (\d+), sender "([^"]*)", approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), app-args "([^"]*)", foreign-apps "([^"]*)", foreign-assets "([^"]*)", app-accounts "([^"]*)", extra-pages (\d+), fee (\d+), first-valid (\d+), last-valid (\d+), genesis-hash "([^"]*)"$`, iBuildAnApplicationTransactionUnit)
	s.Step(`^sign the transaction$`, signTheTransaction)
	s.Step(`^the base(\d+) encoded signed transaction should equal "([^"]*)"$`, theBaseEncodedSignedTransactionShouldEqual)

	//@unit.applications
	s.Step(`^we make a GetAssetByID call for assetID (\d+)$`, weMakeAGetAssetByIDCall)
	s.Step(`^we make a GetApplicationByID call for applicationID (\d+)$`, weMakeAGetApplicationByIDCall)
	s.Step(`^we make a SearchForApplications call with applicationID (\d+)$`, weMakeASearchForApplicationsCall)
	s.Step(`^we make a LookupApplications call with applicationID (\d+)$`, weMakeALookupApplicationsCall)
}
