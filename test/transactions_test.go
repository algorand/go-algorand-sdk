package test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/cucumber/godog"

	"golang.org/x/crypto/ed25519"
)

var signingAccount crypto.Account
var sk1 ed25519.PrivateKey
var addr1 types.Address

func aSigningAccountWithAddressAndMnemonic(address, mnem string) error {
	var err error
	addr1, err = types.DecodeAddress(address)
	if err != nil {
		return err
	}

	sk1, err = mnemonic.ToPrivateKey(mnem)
	account = crypto.Account{
		Address:    addr1,
		PrivateKey: sk1,
		PublicKey:  ed25519.PublicKey(addr1[:]),
	}

	return err
}

func suggestedTransactionParametersTxn(fee int, flatFee string, firstValid, LastValid int, genesisHash, genesisId string) error {
	if flatFee != "true" && flatFee != "false" {
		return fmt.Errorf("flatFee must be either 'true' or 'false'")
	}

	genHash, err := base64.StdEncoding.DecodeString(genesisHash)
	if err != nil {
		return err
	}

	sugParams = types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       genesisId,
		GenesisHash:     genHash,
		FirstRoundValid: types.Round(firstValid),
		LastRoundValid:  types.Round(LastValid),
		FlatFee:         flatFee == "true",
	}

	return nil
}

func signTheTransaction() error {
	var err error
	txid, stx, err = crypto.SignTransaction(sk1, tx)
	return err
}

func theBase64EncodedSignedTransactionShouldEqual(golden string) error {
	gold, err := base64.StdEncoding.DecodeString(golden)
	if err != nil {
		return err
	}
	if !bytes.Equal(gold, stx) {
		toPrint := base64.StdEncoding.EncodeToString(stx)
		return fmt.Errorf("Actual signed transaction does not match the expected. Got %s", toPrint)
	}
	return nil
}

func buildKeyregTransaction(sender, nonparticipation string,
	voteFirst, voteLast, keyDilution int,
	votePkB64, selectionPkB64, stateProofPkB64 string) error {

	if voteFirst < 0 || voteLast < 0 || keyDilution < 0 {
		return fmt.Errorf("Integer arguments cannot be negative")
	}

	nonPartValue, err := strconv.ParseBool(nonparticipation)
	if err != nil {
		return fmt.Errorf("Could not parse nonparticipation value: %v", err)
	}


	tx, err = future.MakeKeyRegTxnWithStateProofKey(sender, nil, sugParams, votePkB64, selectionPkB64, stateProofPkB64, uint64(voteFirst), uint64(voteLast), uint64(keyDilution), nonPartValue)
	return err
}

// this function is a legacy step, we should get rid of it in favor one a newer step that works like
// buildKeyregTransaction above
func buildLegacyAppCallTransaction(
	operation string,
	applicationId int,
	sender, approvalProgram, clearProgram string,
	globalBytes, globalInts, localBytes, localInts int,
	appArgs, foreignApps, foreignAssets, appAccounts string,
	fee, firstValid, lastValid int,
	genesisHash string, extraPages int) error {

	if applicationId < 0 || globalBytes < 0 || globalInts < 0 || localBytes < 0 || localInts < 0 || extraPages < 0 || fee < 0 || firstValid < 0 || lastValid < 0 {
		return fmt.Errorf("Integer arguments cannot be negative")
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		return err
	}

	var approvalP []byte
	if approvalProgram != "" {
		approvalP, err = ioutil.ReadFile("features/resources/" + approvalProgram)
		if err != nil {
			return err
		}
	}

	var clearP []byte
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

	gh, err := base64.StdEncoding.DecodeString(genesisHash)
	if err != nil {
		return err
	}

	// this is only kept to keep compatability with old features
	// going forward, use txnSuggestedParams
	sugParams = types.SuggestedParams{
		Fee:             types.MicroAlgos(uint64(fee)),
		GenesisID:       "",
		GenesisHash:     gh,
		FirstRoundValid: types.Round(firstValid),
		LastRoundValid:  types.Round(lastValid),
		FlatFee:         true,
	}

	switch operation {
	case "create":
		tx, err = future.MakeApplicationCreateTxWithExtraPages(false, approvalP, clearP,
			gSchema, lSchema, args, accs, fApp, fAssets,
			sugParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{}, uint32(extraPages))
	case "update":
		tx, err = future.MakeApplicationUpdateTx(uint64(applicationId), args, accs, fApp, fAssets,
			approvalP, clearP,
			sugParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "call":
		tx, err = future.MakeApplicationCallTx(uint64(applicationId), args, accs,
			fApp, fAssets, types.NoOpOC, approvalP, clearP, gSchema, lSchema,
			sugParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "optin":
		tx, err = future.MakeApplicationOptInTx(uint64(applicationId), args, accs, fApp, fAssets,
			sugParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "clear":
		tx, err = future.MakeApplicationClearStateTx(uint64(applicationId), args, accs, fApp, fAssets,
			sugParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "closeout":
		tx, err = future.MakeApplicationCloseOutTx(uint64(applicationId), args, accs, fApp, fAssets,
			sugParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "delete":
		tx, err = future.MakeApplicationDeleteTx(uint64(applicationId), args, accs, fApp, fAssets,
			sugParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	default:
		err = fmt.Errorf("Unknown opperation: %s", operation)
	}
	return err

}

func TransactionsUnitContext(s *godog.Suite) {
	// @unit.transactions
	s.Step(`^a signing account with address "([^"]*)" and mnemonic "([^"]*)"$`, aSigningAccountWithAddressAndMnemonic)
	s.Step(`^sign the transaction$`, signTheTransaction)
	s.Step(`^the base64 encoded signed transaction should equal "([^"]*)"$`, theBase64EncodedSignedTransactionShouldEqual)
	s.Step(`^the decoded transaction should equal the original$`, theDecodedTransactionShouldEqualTheOriginal)

	// @unit.transactions.keyreg
	s.Step(`^I build a keyreg transaction with sender "([^"]*)", nonparticipation "([^"]*)", vote first (\d+), vote last (\d+), key dilution (\d+), vote public key "([^"]*)", selection public key "([^"]*)", and state proof public key "([^"]*)"$`, buildKeyregTransaction)
	s.Step(`^I build an application transaction with operation "([^"]*)", application-id (\d+), sender "([^"]*)", approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), app-args "([^"]*)", foreign-apps "([^"]*)", foreign-assets "([^"]*)", app-accounts "([^"]*)", fee (\d+), first-valid (\d+), last-valid (\d+), genesis-hash "([^"]*)", extra-pages (\d+)$`, buildLegacyAppCallTransaction)
}
